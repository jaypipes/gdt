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

Alternately, include `github.com/jaypipes/gdt` in your Golang dependency
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

The idea behind `gdt` is to allow test authors to **cleanly** and **clearly**
describe a functional test's **assumptions** and **assertions** in a
declarative format.

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
require:
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
this ([`types_test.go`](examples/books/api/types_test.go)):

```go
package api_test

import (
    "github.com/jaypipes/gdt/examples/books/api"
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
([`failure_test.go`](examples/books/api/failure_test.go)):


```go
package api_test

import (
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httptest"
    "os"
    "strings"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/jaypipes/gdt/examples/books/api"
)

var (
    server *httptest.Server
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
    return strings.TrimSuffix(server.URL, "/") + "/" + strings.TrimPrefix(path, "/")
}

// Register an HTTP server fixture that spins up the API service on a
// random port on localhost
var _ = BeforeSuite(func() {
    logger := log.New(os.Stdout, "http: ", log.LstdFlags)
    c := api.NewControllerWithBooks(logger, nil)
    server = httptest.NewServer(c.Router())
})

var _ = AfterSuite(func() {
    server.Close()
})

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
how `gdt` allows the test author to describe the same assertions
([`failures.yaml`](examples/books/tests/api/failures.yaml)):

```yaml
require:
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
more verbose and cumbersome the Golang test code tends to become. First and
foremost, tests should be *readable*. If they are not readable, then the test's
assertions are not *understandable*. And tests that cannot easily be understood
are often the source of bit rot and technical debt. Worse, tests that aren't
understandable stand a greater chance of having an improper assertion go
undiscovered, leading to tests that validate the wrong behaviour or don't
validate the correct behaviour.

Consider a Ginkgo test case that checks the following behaviour:

* When a book is created via a call to `POST /books`, we are able to get book
 information from the link returned in the HTTP response's `Location` header
* The newly-created book's author name should be set to a known value
* The newly-created book's ID field is a valid UUID
* The newly-created book's publisher has an address containing a known state code

A typical implementation of a Ginkgo Golang test might look like this ([`create_then_get_test.go`](examples/books/api/create_then_get_test.go)):

```go
package api_test

import (
    "bytes"
    "encoding/json"
    "net/http"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/jaypipes/gdt/examples/books/api"
)

var _ = Describe("Books API - POST /books -> GET /books from Location", func() {

    var err error
    var resp *http.Response
    var locURL string
    var authorID, publisherID string

    Describe("proper HTTP GET after POST", func() {

        Context("when creating a single book resource", func() {
            It("should be retrievable via GET {location header}", func() {
                // See https://github.com/onsi/ginkgo/issues/457 for why this
                // needs to be here instead of in the outer Describe block.
                authorID = getAuthorByName("Ernest Hemingway").ID
                publisherID = getPublisherByName("Charles Scribner's Sons").ID
                req := api.CreateBookRequest{
                    Title:       "For Whom The Bell Tolls",
                    AuthorID:    authorID,
                    PublisherID: publisherID,
                    PublishedOn: "1940-10-21",
                    Pages:       480,
                }
                var payload []byte
                payload, err = json.Marshal(&req)
                if err != nil {
                    Fail("Failed to serialize JSON in setup")
                }
                resp, err = http.Post(apiPath("/books"), "application/json", bytes.NewBuffer(payload))
                Ω(err).Should(BeNil())

                // See https://github.com/onsi/ginkgo/issues/70 for why this
                // has to be one giant It() block. The GET tests rely on the
                // result of an earlier POST response (for the Location header)
                // and therefore all of the assertions below much be in a
                // single It() block. :(

                Ω(resp.StatusCode).Should(Equal(201))
                Ω(resp.Header).Should(HaveKey("Location"))

                locURL = resp.Header["Location"][0]

                resp, err = http.Get(apiPath(locURL))
                Ω(err).Should(BeNil())

                Ω(resp.StatusCode).Should(Equal(200))

                var book api.Book

                err := json.Unmarshal([]byte(respJSON(resp)), &book)
                Ω(err).Should(BeNil())

                Ω(IsValidUUID4(book.ID)).Should(BeTrue())
                Ω(book.Author).ShouldNot(BeNil())
                Ω(book.Author.Name).Should(Equal("Ernest Hemingway"))
                Ω(book.Publisher).ShouldNot(BeNil())
                Ω(book.Publisher.Address).ShouldNot(BeNil())
                Ω(book.Publisher.Address.State).Should(Equal("NY"))
            })
        })
    })
})
```

Compare the above test code to the following YAML document that a `gdt` user
might create to describe the same assertions 
([`create_then_get.yaml`](examples/books/tests/api/create_then_get.yaml)):

```yaml
require:
 - books_api
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
         $.id: uuid4
```

## `gdt` test file structure

All gdt test files contain YAML. All `gdt` test files, regardless of their type
(see below), have the following attributes:

* `name`: (optional) string describing the contents of the test file. If
  missing or empty, the filename is used as the name
* `description`: (optional) string with longer description of the test file
  contents
* `type`: (optional) string indicating the type of tests contained in the file.
  `gdt` looks up a test file parser that understands this type of test.
  Defaults to "http"
* `requires`: (optional) list of strings indicating fixtures that will be
  started before any of the tests in the file are run

Depending on the `type` of the test, a parser is invoked to interpret the test
file according to that particular type of test. See the documentation for the
[`gdt.http`](http/README.md) test type for an example of how different types of
tests are handled by an extensible parsing system in `gdt`.
