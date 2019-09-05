package api_test

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/jaypipes/gdt/http"

	"github.com/jaypipes/gdt/examples/books/api"
)

func TestBooksAPI(t *testing.T) {
	// Register an HTTP server fixture that spins up the API service on a
	// random port on localhost
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	c := api.NewControllerWithBooks(logger, nil)
	f := http.NewHTTPServerFixture(c.Router())
	gdt.Fixtures.Register("books_api", f)

	// Construct a new gdt.Runnable from the directory of this file
	_, filename, _, _ := runtime.Caller(0)
	cwd := filepath.Dir(filename)

	ts, err := gdt.From(cwd)
	if err != nil {
		t.Fatal(err)
	}
	ts.Run(t)
}
