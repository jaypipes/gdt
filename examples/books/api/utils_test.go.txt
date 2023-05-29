package api_test

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/jaypipes/gdt/examples/books/api"
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

func IsValidUUID4(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func getAuthorByName(name string) *api.Author {
	for _, book := range data {
		if book.Author.Name == name {
			return book.Author
		}
	}
	return nil
}

func getPublisherByName(name string) *api.Publisher {
	for _, book := range data {
		if book.Publisher.Name == name {
			return book.Publisher
		}
	}
	return nil
}
