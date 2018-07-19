package main

import (
	"encoding/json"
	"golang.org/x/text/language"
	"golang.org/x/text/search"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Record struct {
	Name string `json:"name"`
	Code string `json:"code"`
	URL string `json:"url"`
}

var database []Record
var matcher *search.Matcher

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	if _, ok := r.Form["q"]; !ok {
		json.NewEncoder(w).Encode(database)
		return
	}

	query := r.Form.Get("q")

	if query == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin:", "*")

	var results []Record
	for _, record := range database {
		pattern := matcher.CompileString(query)
		if i, _ := pattern.IndexString(record.Name); i > -1 {
			results = append(results, record)
		}
	}
	json.NewEncoder(w).Encode(results)

	w.Header().Set("ETag", "")

	log.Println(r.Method, r.URL.Path, r.Form.Encode(), time.Since(start))
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

	log.Println("Listening...")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), http.HandlerFunc(handler)))
}
