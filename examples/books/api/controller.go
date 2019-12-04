package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

func errAuthorNotFound(authorID string) error {
	return fmt.Errorf("No such author: %s", authorID)
}

func errPublisherNotFound(publisherID string) error {
	return fmt.Errorf("No such publisher: %s", publisherID)
}

type controller struct {
	lock       *sync.Mutex
	logger     *log.Logger
	authors    map[string]*Author
	publishers map[string]*Publisher
	books      map[string]*Book
}

func NewController(logger *log.Logger) *controller {
	return &controller{
		lock:       &sync.Mutex{},
		logger:     logger,
		authors:    map[string]*Author{},
		publishers: map[string]*Publisher{},
		books:      map[string]*Book{},
	}
}

func NewControllerWithBooks(logger *log.Logger, data []*Book) *controller {
	authors := make(map[string]*Author, 0)
	publishers := make(map[string]*Publisher, 0)
	books := make(map[string]*Book, len(data))
	for _, book := range data {
		if book.Author != nil {
			if book.Author.ID != "" {
				if _, found := authors[book.Author.ID]; !found {
					authors[book.Author.ID] = book.Author
				}
			}
		}
		if book.Publisher != nil {
			if book.Publisher.ID != "" {
				if _, found := publishers[book.Publisher.ID]; !found {
					publishers[book.Publisher.ID] = book.Publisher
				}
			}
		}
		books[book.ID] = book
	}
	return &controller{
		lock:       &sync.Mutex{},
		logger:     logger,
		authors:    authors,
		publishers: publishers,
		books:      books,
	}
}

func (c *controller) Router() http.Handler {
	router := http.NewServeMux()
	router.Handle("/books/", handleBook(c))
	router.Handle("/books", handleBooks(c))
	return router
}

func handleBooks(c *controller) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			postBooks(c, w, r)
			return
		case "PUT":
			putBooks(c, w, r)
			return
		case "GET":
			listBooks(c, w, r)
			return
		}
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	})
}

func handleBook(c *controller) http.Handler {
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
	c *controller,
	w http.ResponseWriter,
	r *http.Request,
	bookID string,
) {
	book := c.getBook(bookID)
	if book == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

func listBooks(c *controller, w http.ResponseWriter, r *http.Request) {
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
	var lbr ListBooksResponse
	lbr.Books = c.listBooks()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&lbr)
}

func postBooks(c *controller, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var cbr CreateBookRequest
	err := decoder.Decode(&cbr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	createdID, err := c.createBook(&cbr)
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

func (c *controller) createBook(cbr *CreateBookRequest) (string, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	author, found := c.authors[cbr.AuthorID]
	if !found {
		return "", errAuthorNotFound(cbr.AuthorID)
	}

	publisher, found := c.publishers[cbr.PublisherID]
	if !found {
		return "", errPublisherNotFound(cbr.PublisherID)
	}
	createdID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	c.books[cbr.Title] = &Book{
		Title:       cbr.Title,
		PublishedOn: cbr.PublishedOn,
		Pages:       cbr.Pages,
		ID:          createdID.String(),
		Author:      author,
		Publisher:   publisher,
	}
	return createdID.String(), nil
}

// putBooks accepts an array of Book entries and creates/replaces the Book
// entries in the API server. Not a great REST API design, but it allows us to
// test the PUT method and array pre-processing for the HTTP test hander
func putBooks(c *controller, w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var rbr ReplaceBooksRequest
	err := decoder.Decode(&rbr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	for _, entry := range rbr {
		// Was ID field set? If so, replace the existing entry, otherwise
		// create a new book
		if entry.ID != "" {
			fmt.Printf("replacing book with ID %s\n", entry.ID)
		} else {
			cbr := CreateBookRequest{
				Title:       entry.Title,
				AuthorID:    entry.AuthorID,
				PublisherID: entry.PublisherID,
				Pages:       entry.Pages,
			}
			_, err := c.createBook(&cbr)
			if err != nil {
				fmt.Printf("XXXXXX: %s", err)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(400)
				if err := json.NewEncoder(w).Encode(err); err != nil {
					panic(err)
				}
				return
			}
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (c *controller) listBooks() []*Book {
	res := make([]*Book, 0, len(c.books))
	for _, book := range c.books {
		res = append(res, book)
	}
	return res
}

func (c *controller) getBook(bookID string) *Book {
	for _, book := range c.books {
		if book.ID == bookID {
			return book
		}
	}
	return nil
}

func (c *controller) Log(args ...interface{}) {
	c.logger.Println(args...)
}

func (c *controller) Panic(fs string, args ...interface{}) {
	c.logger.Fatalf(fs, args...)
}
