package api_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
