package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"../api"
)

var (
	listenAddr string
)

func main() {
	flag.StringVar(&listenAddr, "listen", ":8081", "Books API server listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	c := api.NewController(logger)
	router := http.NewServeMux()
	router.Handle("/books", books(c))

	server := &http.Server{
		Addr:     listenAddr,
		Handler:  logging(logger)(router),
		ErrorLog: logger,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		c.Log("books API server in shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			c.Panic("failed to gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	c.Log("books API server ready for connections. listening on ", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		c.Panic("failed to listen on %s: %v\n", listenAddr, err)
	}

	<-done
	c.Log("books API server stopped")
}

func books(c *api.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			postBooks(c, w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	})
}

func postBooks(c *api.Controller, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/books" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var cbr api.CreateBookRequest
	err := decoder.Decode(&cbr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	createdID, err := c.CreateBook(&cbr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	locURL := fmt.Sprintf("/books/%s", createdID)
	w.Header().Set("Location", locURL)
	w.WriteHeader(http.StatusCreated)
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
