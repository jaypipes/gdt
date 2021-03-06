package api_test

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/jaypipes/gdt/fixtures"
	"github.com/jaypipes/gdt/http"

	"github.com/jaypipes/gdt/examples/books/api"
)

var (
	data struct {
		Authors    interface{}
		Publishers interface{}
		Books      []*api.Book
	}
)

func TestBooksAPI_HTTP(t *testing.T) {
	// Register an HTTP server fixture that spins up the API service on a
	// random port on localhost
	dataFilepath := "testdata/fixtures.json"

	dataFile, err := os.Open(dataFilepath)
	if err != nil {
		panic(err)
	}
	if err = json.NewDecoder(dataFile).Decode(&data); err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "books_api_http: ", log.LstdFlags)
	c := api.NewControllerWithBooks(logger, data.Books)
	apiFixture := http.NewServerFixture(c.Router(), false /* useTLS */)
	gdt.Fixtures.Register("books_api", apiFixture)

	dataFile.Seek(0, io.SeekStart)
	dataFixture, err := fixtures.NewJSONFixture(dataFile)
	if err != nil {
		panic(err)
	}
	gdt.Fixtures.Register("books_data", dataFixture)

	// Construct a new gdt.Runnable from the directory of this file
	_, filename, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(filename)

	ts, err := gdt.From(cwd)
	if err != nil {
		t.Fatal(err)
	}
	ts.Run(t)
}

func TestBooksAPI_HTTPS(t *testing.T) {
	// Register an HTTPS server fixture that spins up the API service on a
	// random port on localhost and a well-known cert for localhost/127.0.0.1
	dataFilepath := "testdata/fixtures.json"

	dataFile, err := os.Open(dataFilepath)
	if err != nil {
		panic(err)
	}
	if err = json.NewDecoder(dataFile).Decode(&data); err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "books_api_http: ", log.LstdFlags)
	c := api.NewControllerWithBooks(logger, data.Books)
	apiFixture := http.NewServerFixture(c.Router(), true /* useTLS */)
	gdt.Fixtures.Register("books_api", apiFixture)

	dataFile.Seek(0, io.SeekStart)
	dataFixture, err := fixtures.NewJSONFixture(dataFile)
	if err != nil {
		panic(err)
	}
	gdt.Fixtures.Register("books_data", dataFixture)

	// Construct a new gdt.Runnable from the directory of this file
	_, filename, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(filename)

	ts, err := gdt.From(cwd)
	if err != nil {
		t.Fatal(err)
	}
	ts.Run(t)
}
