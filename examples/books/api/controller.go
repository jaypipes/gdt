package api

import (
	"fmt"
	"log"
	"sync"
)

func errAuthorNotFound(authorID string) error {
	return fmt.Errorf("No such author: %s", authorID)
}

func errPublisherNotFound(publisherID string) error {
	return fmt.Errorf("No such publisher: %s", publisherID)
}

type Controller struct {
	lock       *sync.Mutex
	logger     *log.Logger
	authors    map[string]*Author
	publishers map[string]*Publisher
	books      map[string]*Book
}

func NewController(logger *log.Logger) *Controller {
	return &Controller{
		lock:       &sync.Mutex{},
		logger:     logger,
		authors:    map[string]*Author{},
		publishers: map[string]*Publisher{},
		books:      map[string]*Book{},
	}
}

func NewControllerWithBooks(logger *log.Logger, data []*Book) *Controller {
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
	return &Controller{
		lock:       &sync.Mutex{},
		logger:     logger,
		authors:    authors,
		publishers: publishers,
		books:      books,
	}
}

func (c *Controller) CreateBook(cbr *CreateBookRequest) (string, error) {
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

	createdID := "newID"

	c.books[cbr.Title] = &Book{
		Title:       cbr.Title,
		PublishedOn: cbr.PublishedOn,
		Author:      author,
		Publisher:   publisher,
	}
	return createdID, nil
}

func (c *Controller) ListBooks() []*Book {
	res := make([]*Book, 0, len(c.books))
	for _, book := range c.books {
		res = append(res, book)
	}
	return res
}

func (c *Controller) GetBook(bookID string) *Book {
	return c.books[bookID]
}

func (c *Controller) Log(args ...interface{}) {
	c.logger.Println(args...)
}

func (c *Controller) Panic(fs string, args ...interface{}) {
	c.logger.Fatalf(fs, args...)
}
