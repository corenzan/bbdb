package main

import (
	"bytes"
	"crypto/md5"
	"encoding/csv"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

type dataRecord struct {
	Name string `json:"name"`
	Code string `json:"code"`
	ISPB string `json:"ispb"`
}

type bufRespWriter struct {
	http.ResponseWriter
	status int
	buffer []byte
}

func (w *bufRespWriter) WriteHeader(status int) {
	w.status = status
}

func (w *bufRespWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.buffer = append(w.buffer, b...)
	return len(b), nil
}

func (w *bufRespWriter) Flush() error {
	if w.status == 0 {
		w.status = 200
	}
	_, err := w.ResponseWriter.Write(w.buffer)
	return err
}

var (
	numericExpr = regexp.MustCompile(`^[0-9]+$`)
	boolExpr    = regexp.MustCompile(`^(y|yes|t|true|1)$`)
	forExpr     = regexp.MustCompile(`(?i)(?:for=)([^(;|,| )]+)`)
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

func v1APIHandler(database []*dataRecord, matcher *search.Matcher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
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

func cachingHandler(ttl time.Duration, etag string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d", ttl/time.Second))
		w.Header().Set("ETag", etag)

		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func remoteAddrHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
			s := strings.Index(fwd, ", ")
			if s == -1 {
				s = len(fwd)
			}
			r.RemoteAddr = fwd[:s]
		} else if fwd := r.Header.Get("X-Real-IP"); fwd != "" {
			r.RemoteAddr = fwd
		} else if fwd := r.Header.Get("Forwarded"); fwd != "" {
			if match := forExpr.FindStringSubmatch(fwd); len(match) > 1 {
				r.RemoteAddr = strings.Trim(match[1], `"`)
			}
		}

		next.ServeHTTP(w, r)
	})
}

func bufRespLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b := &bufRespWriter{w, 0, []byte{}}
		t := time.Now()

		next.ServeHTTP(b, r)
		b.Flush()

		log.Println(r.Method, r.URL.Path, r.Form.Encode(), r.RemoteAddr, b.status, len(b.buffer), time.Since(t))
	})
}

func main() {
	var (
		listenAddr string
		dataSrc    string
	)

	flag.StringVar(&listenAddr, "addr", ":8080", "Address the server will bind to")
	flag.StringVar(&dataSrc, "src", "data.csv", "Path to CSV data source")

	flag.Parse()

	database, etag := loadAndHashData(dataSrc)
	matcher := search.New(language.BrazilianPortuguese, search.Loose)

	handler := v1APIHandler(database, matcher)

	handler = cachingHandler(time.Hour*24*365, etag, handler)
	handler = remoteAddrHandler(handler)
	handler = bufRespLogHandler(handler)

	http.ListenAndServe(listenAddr, handler)
}
