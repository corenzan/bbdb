package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Middleware ...
type Middleware func(http.Handler) http.Handler

// Web ...
type Web struct {
	middlewares []Middleware
}

// Use ...
func (w *Web) Use(mw Middleware) {
	w.middlewares = append(w.middlewares, mw)
}

// New ...
func New() *Web {
	return &Web{
		middlewares: []Middleware{},
	}
}

// Listen ...
func (w *Web) Listen(addr string) {
	handler := http.NotFoundHandler()
	for i := len(w.middlewares) - 1; i >= 0; i-- {
		handler = w.middlewares[i](handler)
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	log.Println("Server is listening to", addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", addr, err)
	}

	<-done
	log.Println("Server stopped")
}
