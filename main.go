package main

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/corenzan/bbdb/web"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

type dataRecord struct {
	Name string `json:"name"`
	Code string `json:"code"`
	ISPB string `json:"ispb"`
}

var (
	numericExpr = regexp.MustCompile(`^[0-9]+$`)
	boolExpr    = regexp.MustCompile(`^(y|yes|t|true|1)$`)
)

func loadAndHashData(path string) ([]*dataRecord, string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip headers
	if _, err = reader.Read(); err != nil {
		panic(err)
	}

	var records []*dataRecord

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		r := &dataRecord{
			ISPB: data[0],
			Name: strings.TrimSpace(data[5]),
		}
		if numericExpr.MatchString(data[2]) {
			r.Code = data[2]
		}
		records = append(records, r)
	}

	hasher := md5.New()

	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(records)
	hasher.Write(b.Bytes())

	return records, hex.EncodeToString(hasher.Sum(nil))
}

func apiHandler(database []*dataRecord, matcher *search.Matcher) web.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				next.ServeHTTP(w, r)
				return
			}

			if err := r.ParseForm(); err != nil {
				log.Print("ParseForm():", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Access-Control-Allow-Origin", "*")

			results := []*dataRecord{}

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

			if boolExpr.MatchString(r.FormValue("code")) {
				tmp := results[:0]
				for _, record := range results {
					if record.Code != "" {
						tmp = append(tmp, record)
					}
				}
				results = tmp
			}

			json.NewEncoder(w).Encode(results)
		})
	}
}

func main() {
	var (
		addr string
		src  string
	)

	flag.StringVar(&addr, "addr", ":8080", "Address the server will listen to")
	flag.StringVar(&src, "src", "data.csv", "Path to CSV data source")

	flag.Parse()

	database, etag := loadAndHashData(src)
	matcher := search.New(language.BrazilianPortuguese, search.Loose)

	w := web.New()

	w.Use(web.LoggingHandler())
	w.Use(web.CachingHandler(time.Hour*24*365, etag))
	w.Use(web.RemoteAddrHandler())
	w.Use(apiHandler(database, matcher))

	w.Listen(addr)
}
