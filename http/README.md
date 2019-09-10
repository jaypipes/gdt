# `gdt/http` - Golang Declarative Testing for HTTP APIs

`gdt/http` is a module for the
[`github.com/jaypipes/gdt`](https://github.com/jaypipes/gdt) testing library
that allows test authors to cleanly describe functional tests of HTTP APIs
using a simple, clear YAML format. `gdt/http` parses YAML files that describe
HTTP requests and assertions about what the HTTP response should contain.

## Installation

`gdt/http` is a Golang library and is intended to be included in your own Golang
application's test code as a Golang package dependency.

Install `gdt/http` into your `$GOPATH` by executing:

```
go get -u github.com/jaypipes/gdt/http
```

## `gdt/http` file format

[`examples/books/tests/api/failures.yaml`](../examples/books/tests/api/failures.yaml):

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

[`examples/books/tests/api/create_then_get.yaml`](../examples/books/tests/api/create_then_get.yaml):

```yaml
require:
 - books_api
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
         $.publisher.address.state: NY
       path_formats:
         $.id: uuid4
```

### `$LOCATION`

TODO(jaypipes)

### Response assertions

TODO(jaypipes)
#### Checking for a string in response body

TODO(jaypipes)
#### Checking for an HTTP header

TODO(jaypipes)
#### Checking for JSON in response

TODO(jaypipes)
