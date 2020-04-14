package main

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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

var boolExpr = regexp.MustCompile("^(on|yes|true|1)$")
var matcher = search.New(language.BrazilianPortuguese, search.Loose)
var database = []*Record{}

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

	results := []*Record{}

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
	if boolExpr.MatchString(compact) {
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

func loadDataSrc(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		database = append(database, &Record{
			Code: record[0],
			Name: record[1],
		})

		if len(record) > 2 {
			if u, err := url.Parse(record[2]); err == nil {
				u.Scheme = "http"
				database[len(database)-1].URL = u.String()
			}
		}
	}
}

func main() {
	var (
		addr string
		src  string
	)

	flag.StringVar(&addr, "addr", ":8080", "Address the server will bind to")
	flag.StringVar(&src, "src", "database.csv", "Path to CSV data source")
	flag.Parse()

	loadDataSrc(src)

	log.Fatal(http.ListenAndServe(addr, http.HandlerFunc(middleware)))
}
