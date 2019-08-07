# `gdt` - The Golang declarative testing framework

`gdt` is a wrapper around the [`gingko`](http://onsi.github.io/ginkgo/) Golang
behavior-driven development (BDD) testing library that allows test authors to
cleanly describe tests in a YAML file. `gdt` reads YAML files that describe a
test and automatically creates the Ginkgo objects in a test suite.

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
above with how `gdt` would allow the test author to validate the same
assertions (`examples/books/failures.yaml`):

```yaml
fixtures:
 - BooksAPI
tests:
 - name: no such book was found
   get: /books/nosuchbook
   response:
     length: 0
     status: 404
 - name: invalid query parameter is supplied
   get: /books?invalidparam=1
   response:
     length: 0
     status: 400
     strings:
       - invalid parameter
```

No more closures and boilerplate function code getting in the way of expressing
the assertions, which should be the focus of the test.
