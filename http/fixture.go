package http

import (
	nethttp "net/http"
	"net/http/httptest"

	"github.com/jaypipes/gdt/interfaces"
)

type httpServerFixture struct {
	handler nethttp.Handler
	server  *httptest.Server
}

func (f *httpServerFixture) Start() {
	f.server = httptest.NewServer(f.handler)
}

func (f *httpServerFixture) Stop() {
	f.server.Close()
}

func (f *httpServerFixture) HasState(key string) bool {
	if key == "http.base_url" {
		return true
	}
	return false
}

func (f *httpServerFixture) State(key string) string {
	if key == "http.base_url" {
		return f.server.URL
	}
	return ""
}

// NewHTTPServerFixture returns a fixture that will start and stop a supplied
// http.Handler. The returned fixture exposes an "http.base_url" state key that
// test cases of type "http" examine to determine the base URL the tests should
// hit
func NewHTTPServerFixture(h nethttp.Handler) interfaces.Fixture {
	return &httpServerFixture{handler: h}
}
