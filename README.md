# `gdt` - The Go Declarative Testing framework ![go test workflow](https://github.com/jaypipes/gdt/actions/workflows/gate-tests.yml/badge.svg)

`gdt` is a testing library that allows test authors to cleanly describe tests
in a YAML file. `gdt` reads YAML files that describe a test's assertions and
then builds a set of Golang structures that the standard Golang
[`testing`](https://golang.org/pkg/testing/) package can execute.

## Introduction

Writing functional tests in Golang can be overly verbose and tedious. When the
code that tests some part of an application is verbose or tedious, then it
becomes difficult to read the tests and quickly understand the assertions the
test is making.

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
coverage of functional tests.

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
       len: 0
     status: 404
```

## Coming from Ginkgo

When using Ginkgo, developers create tests for a particular module (say, the
`books` module) by creating a `books_test.go` file and calling some Ginkgo
functions in a BDD test style. A sample Ginkgo test might look something like
this ([`types_test.go`](examples/books/api/types_test.go.txt)):

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


This is perfectly fine for simple unit tests of Golang code. However, once the
tests begin to call multiple APIs or packages, the Ginkgo Golang tests start to
get cumbersome. Consider the following example of *functionally* testing the
failure modes for a simple HTTP REST API endpoint
([`failure_test.go`](https://github.com/jaypipes/gdt-examples/blob/main/http/api/failure_test.go)):


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

    "github.com/jaypipes/gdt-examples/http/api"
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
([`failures.yaml`](https://github.com/jaypipes/gdt-examples/blob/main/http/tests/api/failures.yaml)):

```yaml
require:
 - books_api
tests:
 - name: no such book was found
   GET: /books/nosuchbook
   response:
     json:
       len: 0
     status: 404
 - name: invalid query parameter is supplied
   GET: /books?invalidparam=1
   response:
     json:
       len: 0
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

A typical implementation of a Ginkgo Golang test might look like this
([`create_then_get_test.go`](https://github.com/jaypipes/gdt-examples/blob/main/http/api/create_then_get_test.go)):

```go
package api_test

import (
    "bytes"
    "encoding/json"
    "net/http"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "github.com/jaypipes/gdt-examples/http/api"
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
                // and therefore all of the assertions below must be in a
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
([`create_then_get.yaml`](https://github.com/jaypipes/gdt-examples/blob/main/http/tests/api/create_then_get.yaml)):

```yaml
require:
 - books_api
 - books_data
tests:
 - name: create a new book
   POST: /books
   data:
     title: For Whom The Bell Tolls
     published_on: 1940-10-21
     pages: 480
     author_id: $.authors.by_name["Ernest Hemingway"].id
     publisher_id: $.publishers.by_name["Charles Scribner's Sons"].id
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

## `gdt` test scenario structure

A `gdt` test scenario (or just "scenario") is simply a YAML file.

All `gdt` scenarios have the following fields:

* `name`: (optional) string describing the contents of the test file. If
  missing or empty, the filename is used as the name
* `description`: (optional) string with longer description of the test file
  contents
* `defaults`: (optional) is a map, keyed by a plugin name, of default options
  and configuration values for that plugin.
* `require`: (optional) list of strings indicating named fixtures that will be
  started before any of the tests in the file are run
* `tests`: list of [`Spec`][basespec] specializations that represent the
  runnable test units in the test scenario.

[basespec]: https://github.com/jaypipes/gdt-core/blob/e1d23e0974447de0bcd273f151edebeebc2b96c6/spec/spec.go#L27-L39

The scenario's `tests` field is the most important and the [`Spec`][basespec]
objects that it contains are the meat of a test scenario.

## `gdt` test spec structure

A spec represents a single *action* that is taken and zero or more
*assertions* that represent what you expect to see resulting from that action.

Each spec is a specialized class of the base [`Spec`][basespec] that deals with
a particular type of test. For example, there is a `Spec` class called `exec`
that allows you to execute arbitrary commands and assert expected result codes
and output. There is a `Spec` class called `http` that allows you to call an
HTTP URL and assert that the response looks like what you expect. Depending on
how you define your test units, `gdt` will parse the YAML definition into one
of these specialized `Spec` classes.

The base `Spec` class has the following fields (and thus all `Spec` specialized
classes inherit these fields):

* `name`: (optional) string describing the test unit.
* `description`: (optional) string with longer description of the test unit.
* `timeout`: (optional) an object containing [timeout information][timeout] for the test
  unit.
* `timeout.after`: a string duration of time the test unit is expected to
  complete within.
* `timeout.expected`: a bool indicating that the test unit is expected to not
  complete before `timeout.after`. This is really only useful in unit testing.
* `wait` (optional) an object containing [wait information][wait] for the test
  unit.
* `wait.before`: a string duration of time that gdt should wait before
  executing the test unit's action.
* `wait.after`: a string duration of time that gdt should wait after executing
  the test unit's action.

[timeout]: https://github.com/jaypipes/gdt-core/blob/e1d23e0974447de0bcd273f151edebeebc2b96c6/types/timeout.go#L11-L22
[wait]: https://github.com/jaypipes/gdt-core/blob/a204ffb055b284a1e943e1f7c7ec578408071d52/types/wait.go

### `exec` test spec structure

An exec spec is a specialization of the base [`Spec`][basespec] that allows
test authors to execute arbitrary commands and assert that the command results
in an expected result code or output.

The [exec `Spec`][execspec] class has the following fields (in addition to all
the base `Spec` fields listed above):

* `exec`: a string with the exact command to execute. You may execute more than
  one command but must include the `shell` field to indicate that the command
  should be run in a shell. It is best practice, however, to simply use
  multiple `exec` specs instead of executing multiple commands in a single
  shell call.
* `shell`: (optional) a string with the specific shell to use in executing the
  command. If empty (the default), no shell is used to execute the command and
  instead the operating system's `exec` family of calls is used.
* `exit_code`: (optional) an integer with the expected exit code from the
  executed command. The default successful exit code is 0 and therefore you do
  not need to specify this if you expect a successful exit code.
* `out`: (optional) a [`PipeAssertions`][pipeassertions] object containing
  assertions about content in `stdout`.
* `out.is`: (optional) a string with the exact contents of `stdout` you expect
  to get.
* `out.contains`: (optional) a list of one or more strings that *all* must be
  present in `stdout`.
* `out.contains_one_of`: (optional) a list of one or more strings of which *at
  least one* must be present in `stdout`.
* `err`: (optional) a [`PipeAssertions`][pipeassertions] object containing
  assertions about content in `stderr`.
* `err.is`: (optional) a string with the exact contents of `stderr` you expect
  to get.
* `err.contains`: (optional) a list of one or more strings that *all* must be
  present in `stderr`.
* `err.contains_one_of`: (optional) a list of one or more strings of which *at
  least one* must be present in `stderr`.

[execspec]: https://github.com/jaypipes/gdt-core/blob/e1d23e0974447de0bcd273f151edebeebc2b96c6/types/timeout.go#L11-L22
[pipeassertions]: https://github.com/jaypipes/gdt-core/blob/debff55b194b222a0db65d20ffce94f0aba68504/exec/assertions.go#L16-L27

## Contributing and acknowledgements

`gdt` was inspired by [Gabbi](https://github.com/cdent/gabbi), the excellent
Python declarative testing framework. `gdt` tries to bring the same clear,
concise test definitions to the world of Golang functional testing.

Contributions to `gdt` are welcomed! Feel free to open a Github issue or submit
a pull request.
