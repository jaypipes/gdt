# `gdt` - The Golang declarative testing framework

`gdt` is a wrapper around the [`gingko`](http://onsi.github.io/ginkgo/) Golang
testing library that allows test authors to cleanly describe tests in a YAML
file. `gdt` reads YAML files that describe a test and automatically creates the
Ginkgo objects in a test suite.

## Introduction

When using Gingkgo, developers create tests for a particular module (say, the
`books` module) by creating a `books_test.go` file and calling some Ginkgo
functions in a BDD test style. A sample Ginkgo test might look something like
this:

```go
package books_test

import (
    . "/path/to/books"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Book", func() {
    var (
        longBook  Book
        shortBook Book
    )

    BeforeEach(func() {
        longBook = Book{
            Title:  "Les Miserables",
            Author: "Victor Hugo",
            Pages:  1488,
        }

        shortBook = Book{
            Title:  "Fox In Socks",
            Author: "Dr. Seuss",
            Pages:  24,
        }
    })

    Describe("Categorizing book length", func() {
        Context("With more than 300 pages", func() {
            It("should be a novel", func() {
                Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
            })
        })

        Context("With fewer than 300 pages", func() {
            It("should be a short story", func() {
                Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
            })
        })
    })
})
```


This is perfectly great for simple unit tests of Golang code. However, once the
tests begin to call multiple APIs or packages, the Ginkgo Golang tests start to
get cumbersome:


```go
Describe("Books API - GET /books failures", func() {
    var client APIClient
    var booksServer BooksServer
    var response chan APIResponse

    BeforeEach(func() {
        response = make(chan APIResponse, 1)
        booksServer = NewBookServer()
        client = NewAPIClient(booksServer)
    })

    Describe("failure modes", func() {
        AssertZeroJSONLength := func() {
            It("should not include JSON in the response", func() {
                Ω((<-response).JSON).Should(BeZero())
            })
        }

        Context("when no such book was found", func() {
            BeforeEach(func() {
                client.Get("/books/nosuchbook", response)
            })

            AssertZeroJSONLength()

            It("should return 404", func() {
                Ω((<-response).StatusCode).Should(Equal(404))
            })
        })

        Context("when an invalid query parameter is supplied", func() {
            BeforeEach(func() {
                client.Get("/books?invalidparam=1", response)
            })

            AssertZeroJSONLength()

            It("should return 400", func() {
                Ω((<-response).StatusCode).Should(Equal(400))
            })
            It("should indicate invalid query parameter", func() {
                Ω((<-response).Body).Should(Contain("invalid parameter"))
            })
        })
    })
})
```

The above Ginkgo Golang test code obscures what is being tested. Compare the
above with how `gdt` would allow the test author to describe the same
assertions (`examples/books/tests/failures.yaml`):

```yaml
fixtures:
 - BooksAPI
tests:
 - name: no such book was found
   GET: /books/nosuchbook
   response:
     json:
       length: 0
     status: 404
 - name: invalid query parameter is supplied
   GET: /books?invalidparam=1
   response:
     json:
       length: 0
     status: 400
     strings:
       - invalid parameter
```

No more closures and boilerplate function code getting in the way of expressing
the assertions, which should be the focus of the test.

The more intricate the assertions being verified by the test, generally the
more verbose and cumbersome the Golang test code becomes. First and foremost,
tests should be *readable*. If they are not readable, then the test's
assertions are not *understandable*. And tests that cannot easily be understood
are often the source of bit rot and technical debt.

Consider a Ginkgo Golang test case that checks the following behavior:

* When a book is created via a call to `POST /books`, we are able to get book
 information from the link returned in the HTTP response's `Location` header
* The newly-created book's author name should be set to a known value
* The newly-created book's ID field is a valid UUID
* The newly-created book's publisher has an address containing a known state code

A typical implementation of a Ginkgo Golang test might look like this:

```go
package books_test

import (
    "encoding/json"
    "regex"

    . "/path/to/books"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func IsValidUUID(uuid string) bool {
    r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
    return r.MatchString(uuid)
}

Describe("Books API", func() {
    var client APIClient
    var booksServer books.BooksServer
    var response chan APIResponse
    var newBookURL string

    var authorID := getAuthorFixtureByName("Ernest Hemingway").ID
    var publisherID := getPublisherFixtureByName("Charles Scribner's Sons").ID

    BeforeEach(func() {
        response = make(chan APIResponse, 1)
        booksServer = NewBookServer()
        client = NewAPIClient(booksServer)
    })

    Describe("proper HTTP GET after POST", func() {
        Context("when creating a single book resource", func() {
            BeforeEach(func() {
                req := books.CreateBookAPIRequest{
                    Name: "For Whom The Bell Tolls",
                    AuthorID: authorID,
                    PublisherID: publisherID,
                    PublishedOn: "1940-10-21",
                }
                payload, err := json.Marshal()
                if err != nil {
                    Fail("Failed to serialize JSON in setup")
                }
                client.Post("/books", payload)
            })

            It("should return 201", func() {
                Ω((<-response).StatusCode).Should(Equal(200))
            })

            It("should return a Location HTTP header", func() {
                Ω((<-response).Headers).Should(Contain("Location"))
            })

            newBookURL := (<-response).Headers["Location"]
        }

        Context("when creating a single book resource", func() {
            BeforeEach(func() {
                client.Get(newBookURL, response)
            })

            var book books.Book
            err := json.Unmarshal([]byte((<-response).JSON)), &book)

            It("should have valid JSON", func() {
                Ω(err).Should(BeNil())
            })

            It("should have a UUID as ID attribute", func() {
                Ω(IsValidUUID(book.ID)).Should(BeTrue())
            })

            It("should have an Author sub-object", func() {
                Ω(book.Author).ShouldNot(BeNil())
            })

            It("should have an Author Name as expected", func() {
                Ω(book.Author.Name).Should(Equal("Ernest Hemingway"))
            })

            It("should have a Publisher sub-object", func() {
                Ω(book.Publisher).ShouldNot(BeNil())
            })

            It("should have a Publisher.Address sub-object", func() {
                Ω(book.Publisher.Address).ShouldNot(BeNil())
            })

            It("should have an Publisher State as expected", func() {
                Ω(book.Publisher.State).Should(Equal("New York"))
            })
        })
    })
})
```

Compare the above test code to the following YAML document that a `gdt` user
might create to describe the same assertions 
(`examples/books/tests/create-then-get.yaml`):

```yaml
fixtures:
 - BooksAPI
 - Authors
 - Publishers
tests:
 - name: create a new book
   POST: /books
   data:
     name: Ernest Hemingway
     authorID: $FIXTURES['Authors']['Ernest Hemingway']['ID']
     publisherID: $FIXTURES['Publishers']["Charles Scribner's Sons"]['ID']
   response:
     status: 201
     headers:
      - Location
 - name: look up that created book
   GET: $LOCATION
   response:
     status: 200
     json:
       paths:
         $.author.name: Ernest Hemingway
         $.publisher.address.state: New York
       path_formats:
         $.id: uuid
```
