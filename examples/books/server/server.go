package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	server := &http.Server{
		Addr:     listenAddr,
		Handler:  logging(logger)(c.Router()),
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
