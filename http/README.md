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
* `data`: (optional) if present, will be encoded into the HTTP request
  payload. Elements of the `data` structure may be JSONPath expressions (see [below](#use-jsonpath-expressions-to-substitute-fixture-data))
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

* `length`: (optional) integer representing the number of bytes in the
 resulting JSON object after successfully parsing the HTTP response body
* `paths`: (optional) map of strings where the keys of the map are JSONPath
  expressions and the values of the map are the expected value to be found when
  evaluating the JSONPath expression
* `path_formats`: (optional) map of strings where the keys of the map are
  JSONPath expressions and the values of the map are the expected format of the
  value to be found when evaluating the JSONPath expression. See the
  [list of valid format strings](#valid-format-strings)
* `schema`: (optional) string containing a filepath to a JSONSchema document.
  If present, the JSON included in the HTTP response will be validated against
  this JSONSChema document.

### Specify HTTP request payload

The `data` attribute of the test unit is used to specify a payload to be
encoded into the HTTP request body. By default, the contents of the `data`
attribute are encoded as JSON.

TODO(jaypipes): Support non-JSON encoding.

The `data` attribute is especially useful for testing of `POST` and `PUT`
requests, where you want to send data to the server to create or update some
resource.

For example, suppose the `POST /books` URL accepts some JSON-encoded data with
information about the to-be-created book's author, title, publisher, etc.

To test the `POST /books` functionality, a test author might use the following
test unit:

```yaml
 - name: create a new book
   POST: /books
   data:
     title: For Whom The Bell Tolls
     published_on: 1940-10-21
     pages: 480
     author_id: "1"
     publisher_id: "1"
   response:
     status: 201
```

The above `data` would be encoded into the following HTTP request body:

```json
{
     "title": "For Whom The Bell Tolls",
     "published_on": "1940-10-21",
     "pages": 480,
     "author_id": "1",
     "publisher_id": "1"
}
```

#### Use JSONPath expressions to substitute fixture data

Often, you will want to reference some information in a fixture instead of
hard-coding values in the `data` contents.

Consider this test unit:

```yaml
 - name: create a new book
   POST: /books
   data:
     title: For Whom The Bell Tolls
     published_on: 1940-10-21
     pages: 480
     author_id: "1"
     publisher_id: "1"
   response:
     status: 201
```

Hard-coding the string `"1"` for the `author_id` and `publisher_id` values is
fragile and non-descriptive. It requires the reader of the test to know which
author has an ID of "1" and which publisher has an ID of "1".

It would be much more readable if we could replace those hard-coded `"1"`
values with a reference to some fixture data:

```yaml
requires:
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
```

The test reader can now better understand what value is being placed into the
"author_id" field of the HTTP request payload: the ID value of the author whose
name is "Ernest Hemingway".

The JSONPath expressions that replaced the hard-coded `"1"` values are
evaluated by the fixtures associated with a test file. The
A `gdt.fixtures.JSONFixture` fixture is designed to evaluate JSONPath
expressions for the data defined in the fixture.

Assume I have a file `testdata/fixtures.json` that looks like this:

```json
{
    "books": [
        {
            "id": "12ac1b94-5667-461e-80cb-ba8619cae61a",
            "title": "Old Man and the Sea",
            "published_on": "1952-10-01",
            "pages": 127,
            "author": {
                "name": "Ernest Hemingway",
                "id": "1"
            },
            "publisher": {
                "name": "Charles Scribner's Sons",
                "id": "1",
                "address": {
                    "address": "153–157 Fifth Avenue",
                    "city": "New York City",
                    "state": "NY",
                    "postal_code": "10010",
                    "country_code": "US"
                }
            }
        }
    ],
    "authors": {
        "by_name": {
            "Ernest Hemingway": {
                "id": 1
            }
        }
    },
    "publishers": {
        "by_name": {
            "Charles Scribner's Sons": {
                "id": 1
            }
        }
    }
}
```

We can register a `gdt.fixtures.JSONFixture` that contains the data in
`testdata/fixtures.json`:

```go
	dataFilepath := "testdata/fixtures.json"

	dataFile, _ := os.Open(dataFilepath)
	dataFixture, err := fixtures.NewJSONFixture(dataFile)
	if err != nil {
		panic(err)
	}
	gdt.Fixtures.Register("books_data", dataFixture)
```

To reference any of the data in your `gdt.fixtures.JSONFixture` from your test unit, just make sure the fixture is listed in the test file's `requires` field:

```
requires:
 - books_data
```

Then you can grab any data in the fixture using a JSONPath expression in the `data` contents:

```yaml
   data:
     author_id: $.authors.by_name["Ernest Hemingway"].id
```

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

The `url` attribute of an HTTP test spec can be the string `$LOCATION`. When
this is set, the HTTP request will be to the URL specified in the *previous
HTTP response's* Location HTTP header. This is an easy shortcut for testing a
series of ordered HTTP requests, where the first HTTP request (typically a
`POST` or `PUT` to a particular resource) responds with a Location HTTP header
pointing to a URL that can have issued an HTTP `GET` request to return
information about the previously created or mutated resource.

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
         $.publisher.address.state: NY
       path_formats:
         $.id: uuid4
```

[`examples/books/tests/api/get_books.yaml`](../examples/books/tests/api/get_books.yaml):

```yaml
require:
 - books_api
 - books_data
tests:
 - name: list all books
   GET: /books
   response:
     status: 200
     json:
       schema: schemas/get_books.json
```

with the contents of [`examples/books/tests/api/schemas/get_books.json`](../examples/books/tests/api/schemas/get_books.json):

```json
{
  "$id": "/schemas/get_books.json",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "books": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string"
          },
          "title": {
            "type": "string"
          },
          "pages": {
            "type": "number"
          },
          "author": {
            "type": "object",
            "properties": {
              "id": {
                "type": "string"
              },
              "name": {
                "type": "string"
              }
            },
            "require": [ "id", "name" ]
          },
          "publisher": {
            "type": "object",
            "properties": {
              "id": {
                "type": "string"
              },
              "name": {
                "type": "string"
              }
            },
            "require": [ "id", "name" ]
          }
        },
        "required": [ "id", "title", "author" ]
      }
    }
  },
  "required": [ "books" ]
}
```
