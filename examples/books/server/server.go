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
	"strings"
	"time"

	"github.com/jaypipes/gdt/examples/books/api"
)

var (
	listenAddr   string
	dataFilepath string
)

func main() {
	flag.StringVar(&listenAddr, "listen", ":8081", "Books API server listen address")
	flag.StringVar(&dataFilepath, "data-filepath", "books.json", "File with Books JSON data")
	flag.Parse()

	var data []*api.Book
	dataFile, err := os.Open(dataFilepath)
	if err != nil {
		panic(err)
	}
	if err = json.NewDecoder(dataFile).Decode(&data); err != nil {
		panic(err)
	}

	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	c := api.NewControllerWithBooks(logger, data)
	router := http.NewServeMux()
	router.Handle("/books/", handleBook(c))
	router.Handle("/books", handleBooks(c))

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

func handleBooks(c *api.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			postBooks(c, w, r)
			return
		case "GET":
			listBooks(c, w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	})
}

func handleBook(c *api.Controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			bookID := strings.Replace(r.URL.Path, "/books/", "", 1)
			getBook(c, w, r, string(bookID))
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	})
}

func getBook(
	c *api.Controller,
	w http.ResponseWriter,
	r *http.Request,
	bookID string,
) {
	book := c.GetBook(bookID)
	if book == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func listBooks(c *api.Controller, w http.ResponseWriter, r *http.Request) {
	// Our GET /books endpoint only supports a "sort" parameter
	params := r.URL.Query()
	if len(params) > 0 {
		if _, found := params["sort"]; !found {
			var msg string
			for key := range params {
				msg = fmt.Sprintf("invalid parameter: %s", key)
				break
			}
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
	}
	var lbr api.ListBooksResponse
	lbr.Books = c.ListBooks()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&lbr)
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
