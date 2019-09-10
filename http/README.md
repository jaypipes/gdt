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

## `gdt/http` test file structure

The `gdt/http` test file [parser](parse.go) parses a test file with type
"http". It parses the test file into an object with the following attributes:

* `tests`: list of test unit objects that describe a test of an HTTP request and
  response

Each of the test unit objects have the following attributes:

* `name`: (optional) string describing the individual test. If missing or
  empty, the test unit's name is a string with the request HTTP method and path
* `description`: (optional) string with a longer description of the test unit
* `method`: (optional) string with the HTTP verb to use. Defaults to "GET" if
  `url` attribute is non-empty
* `url`: (optional) string with the path or URL to use for the HTTP request. If
  missing, one of the `GET`, `POST`, `PATCH`, `DELETE` or `PUT` shortcut
  attributes must be non-empty
  * `GET`: (optional) string with the path or URL to issue an HTTP GET request
  * `POST`: (optional) string with the path or URL to issue an HTTP POST request
  * `PUT`: (optional) string with the path or URL to issue an HTTP PUT request
  * `PATCH`: (optional) string with the path or URL to issue an HTTP PATCH request
  * `DELETE`: (optional) string with the path or URL to issue an HTTP DELETE request
* `response`: (optional) object describing the **assertions** to make about the
  HTTP response received after issuing the HTTP request

The `response` object has the following attributes:

* `status`: (optional) integer corresponding to the expected HTTP status code
  of the HTTP response
* `strings`: (optional) list of strings that should appear in the body of the
  HTTP response
* `json`: (optional) object describing the assertions to make about JSON
  content in the HTTP response body

The `json` object has the following attributes:

* `paths`: (optional) map of strings where the keys of the map are JSONPath
  expressions and the values of the map are the expected value to be found when
  evaluating the JSONPath expression
* `path_formats`: (optional) map of strings where the keys of the map are
  JSONPath expressions and the values of the map are the expected format of the
  value to be found when evaluating the JSONPath expression. See the
  [list of valid format strings](#valid-format-strings)

### Specify expected response values (`response.json.paths`)

When you want to validate the structure of the returned JSON object in an HTTP
response body, you use the `response.json.paths` attribute of the test unit.

This attribute is a map of string to string, where the map keys are JSONPath
expressions and the map values are the expected value when evaluating that
JSONPath expression.

For example, let's say you want to verify that an HTTP `GET` request to the
`/books` URL returns an HTTP response that contains a list of JSON objects, and
that the first JSON object in that list contains a field, "title", that
contains the string "For Whom the Bell Tolls". You would write the test unit
like so:

```yaml
tests:
  - GET: /books
    response:
      json:
        paths:
          - $[0].title: For Whom the Bell Tolls
```

### Specify expected response value format (`response.json.path_formats`)

When you want to validate that a certain field in a returned JSON object from
an HTTP response matches a particular common format, you use the
`response.json.path_formats` attribute of the test unit.

This attribute is a map of string to string, where the map keys are JSONPath
expressions and the map values are the [type of format](#valid-format-strings)
that the value to be found at the JSONPath expression should have.

For example, let's say you want to verify that an HTTP `GET` request to the
`/books/thebook` URL returned an HTTP response that contains a JSON object
having a "id" field, and that the value of that field is a valid version 4
UUID. You would write the test unit like so:

```yaml
tests:
  - GET: /books/thebook
    response:
      json:
        path_formats:
          - $.id: uuid
```

The `$.id` string is a JSONPath expression that selects the value of the field
called "id" from the top-level document/object. The `uuid4` string indicates
the expected format of that value.

#### Valid format strings

The currently supported format strings are:

* "uuid": must be any version of UUID
* "uuid4": must be a UUID version 4

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

## Examples

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
