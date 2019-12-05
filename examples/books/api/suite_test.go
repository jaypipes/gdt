package api_test

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jaypipes/gdt/examples/books/api"
)

var (
	server *httptest.Server
	data   []*api.Book
)

// Register an HTTP server fixture that spins up the API service on a
// random port on localhost
var _ = BeforeSuite(func() {
	dataFilepath := "testdata/books.json"

	dataFile, err := os.Open(dataFilepath)
	if err != nil {
		panic(err)
	}
	if err = json.NewDecoder(dataFile).Decode(&data); err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	c := api.NewControllerWithBooks(logger, data)
	server = httptest.NewServer(c.Router())
})

var _ = AfterSuite(func() {
	server.Close()
})

func TestBooksAPI_Ginkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books API Suite")
}
