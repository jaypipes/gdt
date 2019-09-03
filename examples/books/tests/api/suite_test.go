package api_test

import (
	"log"
	"os"
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/jaypipes/gdt/http"

	"github.com/jaypipes/gdt/examples/books/api"
)

func TestBooksAPI(t *testing.T) {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	c := api.NewControllerWithBooks(logger, nil)
	f := http.NewHTTPServerFixture(c.Router())
	gdt.Fixtures.Register("books_api", f)
	tc, err := gdt.FromFile(t, "failures.yaml")
	if err != nil {
		t.Fatal(err)
	}
	tc.Run()
}
