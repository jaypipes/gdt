# `gdt` - The Golang Declarative Testing framework

`gdt` is a testing library that allows test authors to cleanly describe tests
in a YAML file. `gdt` reads YAML files that describe a test's assertions and
then builds a set of Golang structures that the standard Golang
[`testing`](https://golang.org/pkg/testing/) package can execute.

## Installation

`gdt` is a Golang library and is intended to be included in your own Golang
application's test code as a Golang package dependency.

Install `gdt` into your `$GOPATH` by executing:

```
go get -u github.com/jaypipes/gdt
```

Alternately, include "github.com/jaypipes/gdt" in your Golang dependency
management of choice.

## Introduction

Writing functional tests in Golang can be overly verbose and tedious. When the
code that functionally tests some part of an application is verbose or tedious,
then it becomes difficult to read the tests and quickly understand the
assertions the test is making.

The more difficult it is to understand the test assertions or the test setups
and assumptions, the greater the chance that the test improperly validates the
application behaviour. Furthermore, test code that is cumbersome to read is
prone to bit-rot due to its high maintenance cost. This is particularly true
for code that verifies an application's integration points with *other*
applications via an API.

The idea behind `gdt` is to allow test authors to cleanly and clearly describe
a functional test's assumptions and assertions in a declarative format.
Separating the *description* of a test's assumptions (setup) and assertions
from the Golang code that actually performs the test assertions leads to tests
that are easier to read and understand. This allows developers to spend *more
time writing code* and less time copy/pasting boilerplate test code. Due to the
easier test comprehension, `gdt` also encourages writing greater quality and
coverage of functional tests due to easier test comprehension.

Instead of developers writing code that looks like this:

```go
var _ = Describe("Books API - GET /books failures", func() {
    var response *http.Response
    var err error
    var testPath = "/books/nosuchbook"

    BeforeEach(func() {
        response, err = http.Get(apiPath(testPath))
        Ω(err).Should(BeZero())
    })

    Describe("failure modes", func() {
        Context("when no such book was found", func() {
            It("should not include JSON in the response", func() {
                Ω(respJSON(response)).Should(BeZero())
            })
            It("should return 404", func() {
                Ω(response.StatusCode).Should(Equal(404))
            })
        })
    })
})
```

they can instead have a test that looks like this:


```yaml
setup:
 - books_api
tests:
 - name: no such book was found
   GET: /books/nosuchbook
   response:
     json:
       length: 0
     status: 404
```

## Coming from Ginkgo

When using Ginkgo, developers create tests for a particular module (say, the
`books` module) by creating a `books_test.go` file and calling some Ginkgo
functions in a BDD test style. A sample Ginkgo test might look something like
this ([`examples/books/api/types_test.go`](examples/books/api/types_test.go)):

```go
package api_test

import (
    api "."
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Books API Types", func() {
    var (
        longBook  api.Book
        shortBook api.Book
    )

    BeforeEach(func() {
        longBook = api.Book{
            Title: "Les Miserables",
            Pages: 1488,
            Author: &api.Author{
                Name: "Victor Hugo",
            },
        }

        shortBook = api.Book{
            Title: "Fox In Socks",
            Pages: 24,
            Author: &api.Author{
                Name: "Dr. Seuss",
            },
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
get cumbersome. Consider the following example of *functionally* testing the
failure modes for a simple HTTP REST API endpoint
([`examples/books/api/failure_test.go`](examples/books/api/failure_test.go)):


```go
package api_test

import (
    "io/ioutil"
    "net/http"
    "os"
    "strings"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

const (
    defaultAPIServerURL = "http://localhost:8081"
)

// respJSON returns a string if the supplied HTTP response body is JSON,
// otherwise the empty string
func respJSON(r *http.Response) string {
    if r == nil {
        return ""
    }
    if !strings.HasPrefix(r.Header.Get("content-type"), "application/json") {
        return ""
    }
    bodyStr, _ := ioutil.ReadAll(r.Body)
    return string(bodyStr)
}

// respText returns a string if the supplied HTTP response has a text/plain
// content type and a body, otherwise the empty string
func respText(r *http.Response) string {
    if r == nil {
        return ""
    }
    if !strings.HasPrefix(r.Header.Get("content-type"), "text/plain") {
        return ""
    }
    bodyStr, _ := ioutil.ReadAll(r.Body)
    return string(bodyStr)
}

func apiPath(path string) string {
    serverURL, found := os.LookupEnv("EXAMPLES_BOOKS_API_SERVER_URL")
    if !found {
        serverURL = defaultAPIServerURL
    }
    return strings.TrimSuffix(serverURL, "/") + "/" + strings.TrimPrefix(path, "/")
}

var _ = Describe("Books API - GET /books failures", func() {
    var response *http.Response
    var err error
    var testPath string

    BeforeEach(func() {
        response, err = http.Get(apiPath(testPath))
        Ω(err).Should(BeZero())
    })

    Describe("failure modes", func() {
        AssertZeroJSONLength := func() {
            It("should not include JSON in the response", func() {
                Ω(respJSON(response)).Should(BeZero())
            })
        }

        Context("when no such book was found", func() {
            JustBeforeEach(func() {
                testPath = "/books/nosuchbook"
            })

            AssertZeroJSONLength()

            It("should return 404", func() {
                Ω(response.StatusCode).Should(Equal(404))
            })
        })

        Context("when an invalid query parameter is supplied", func() {
            JustBeforeEach(func() {
                testPath = "/books?invalidparam=1"
            })

            AssertZeroJSONLength()

            It("should return 400", func() {
                Ω(response.StatusCode).Should(Equal(400))
            })
            It("should indicate invalid query parameter", func() {
                Ω(respText(response)).Should(ContainSubstring("invalid parameter"))
            })
        })
    })
})
```

The above test code obscures what is being tested by cluttering the test
assertions with the Golang closures and accessor code. Compare the above with
how `gdt` would allow the test author to describe the same assertions
(`examples/books/api/tests/failures.yaml`):

```yaml
setup:
 - books_api
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
are often the source of bit rot and technical debt. Worse, tests that aren't
understandable stand a greater chance of having an improper assertion go
undiscovered, leading to tests that validate the wrong behaviour or don't
validate the correct behaviour.

Consider a Ginkgo Golang test case that checks the following behaviour:

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
                payload, err := json.Marshal(&req)
                if err != nil {
                    Fail("Failed to serialize JSON in setup")
                }
                client.Post("/books", payload)
            })

            It("should return 201", func() {
                Ω((<-response).StatusCode).Should(Equal(200))
            })

            It("should return a Location HTTP header", func() {
                Ω((<-response).Headers).Should(HaveKey("Location"))
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
setup:
 - books_api
fixtures:
 - authors_by_name
 - publishers_by_name
tests:
 - name: create a new book
   POST: /books
   data:
     title: For Whom The Bell Tolls
     published_on: 1940-10-21
     pages: 480
     author_id: $authors_by_name['Ernest Hemingway']['ID']
     publisher_id: $publishers_by_name["Charles Scribner's Sons"]['ID']
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
