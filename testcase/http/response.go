package api_test

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type response struct {
	*http.Response
}

// JSON returns a string if the supplied HTTP response body is JSON,
// otherwise the empty string
func (r *response) JSON() string {
	if r == nil {
		return ""
	}
	if !strings.HasPrefix(r.Header.Get("content-type"), "application/json") {
		return ""
	}
	bodyStr, _ := ioutil.ReadAll(r.Body)
	return string(bodyStr)
}

// Text returns a string if the supplied HTTP response has a text/plain
// content type and a body, otherwise the empty string
func (r *response) Text() string {
	if r == nil {
		return ""
	}
	if !strings.HasPrefix(r.Header.Get("content-type"), "text/plain") {
		return ""
	}
	bodyStr, _ := ioutil.ReadAll(r.Body)
	return string(bodyStr)
}
