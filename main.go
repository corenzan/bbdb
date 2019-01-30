package main

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

type Record struct {
	Name string `json:"name"`
	Code string `json:"code"`
	URL  string `json:"url"`
}

var boolPat = regexp.MustCompile("^(on|yes|true|1)$")
var database []Record
var matcher *search.Matcher

func hash(i interface{}) [16]byte {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(i)
	return md5.Sum(b.Bytes())
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Print("ParseForm():", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	results := []Record{}

	query := r.FormValue("q")
	if query == "" {
		results = database
	} else {
		for _, record := range database {
			match := matcher.CompileString(query)
			if i, _ := match.IndexString(record.Code + " " + record.Name); i > -1 {
				results = append(results, record)
			}
		}
	}

	compact := r.FormValue("compact")
	if boolPat.MatchString(compact) {
		compacted := results[:0]
		for _, record := range results {
			if record.Code != "" {
				compacted = append(compacted, record)
			}
		}
		results = compacted
	}

	w.Header().Set("ETag", fmt.Sprintf("%x", hash(results)))
	json.NewEncoder(w).Encode(results)
}

func middleware(w http.ResponseWriter, r *http.Request) {
	checkpoint := time.Now()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin:", "*")

	handler(w, r)

	log.Println(r.Method, r.URL.Path, r.Form.Encode(), time.Since(checkpoint))
}

func main() {
	bytes, err := ioutil.ReadFile("database.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal(bytes, &database); err != nil {
		log.Fatal(err)
	}

	matcher = search.New(language.BrazilianPortuguese, search.Loose)

	var addr string
	flag.StringVar(&addr, "addr", "8080", "address the server will bind to")
	flag.Parse()

	log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(middleware)))
}
