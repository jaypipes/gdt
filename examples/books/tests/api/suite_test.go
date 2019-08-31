package api_test

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jaypipes/gdt"
	_ "github.com/jaypipes/gdt/http"

	"github.com/jaypipes/gdt/examples/books/api"
)

type booksAPIFixture struct {
	server *httptest.Server
}

func (f *booksAPIFixture) Start() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	c := api.NewControllerWithBooks(logger, nil)
	f.server = httptest.NewServer(c.Router())
}

func (f *booksAPIFixture) Stop() {
	f.server.Close()
}

func TestBooksAPI(t *testing.T) {
	f := booksAPIFixture{}
	gdt.RegisterFixture(&f, "books_api")
	tc, err := gdt.FromFile(t, "failures.yaml")
	if err != nil {
		t.Fatal(err)
	}
	tc.Run()
}
