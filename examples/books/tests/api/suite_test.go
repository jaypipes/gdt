package api_test

import (
	"testing"

	"github.com/jaypipes/gdt"
	_ "github.com/jaypipes/gdt/http"
)

type booksAPIFixture struct {
}

func (f *booksAPIFixture) Start() {

}
func (f *booksAPIFixture) Cleanup() {

}

func TestBooksAPI(t *testing.T) {
	f := booksAPIFixture{}
	gdt.RegisterFixture(&f, "books_api")
	gdt.FromFile(t, "failures.yaml")
}
