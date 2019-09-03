package api_test

import (
	"fmt"
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
	fmt.Printf("started books_api server at: %s\n", f.server.URL)
}

func (f *booksAPIFixture) Stop() {
	fmt.Println("stopping books_api server")
	f.server.Close()
}

func (f *booksAPIFixture) HasState(key string) bool {
	if key == "http.base_url" {
		return true
	}
	return false
}

func (f *booksAPIFixture) State(key string) string {
	if key == "http.base_url" {
		return f.server.URL
	}
	return ""
}

func TestBooksAPI(t *testing.T) {
	f := booksAPIFixture{}
	gdt.Fixtures.Register("books_api", &f)
	tc, err := gdt.FromFile(t, "failures.yaml")
	if err != nil {
		t.Fatal(err)
	}
	tc.Run()
}
