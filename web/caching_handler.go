package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// CachingHandler ...
func CachingHandler(ttl time.Duration, etag string) Middleware {
	return func(next http.Handler) http.Handler {
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
}
