package http

import (
	nethttp "net/http"
	"net/http/httptest"
	"strings"

	"github.com/jaypipes/gdt"
)

const (
	FIXTURE_STATE_KEY_BASE_URL = "http.base_url"
	FIXTURE_STATE_KEY_CLIENT   = "http.client"
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
	if strings.ToLower(key) == FIXTURE_STATE_KEY_BASE_URL {
		return true
	}
	return false
}

func (f *httpServerFixture) State(key string) interface{} {
	key = strings.ToLower(key)
	switch key {
	case FIXTURE_STATE_KEY_BASE_URL:
		return f.server.URL
	case FIXTURE_STATE_KEY_CLIENT:
		return f.server.Client()
	}
	return ""
}

// NewHTTPServerFixture returns a fixture that will start and stop a supplied
// http.Handler. The returned fixture exposes an "http.base_url" state key that
// test cases of type "http" examine to determine the base URL the tests should
// hit
func NewHTTPServerFixture(h nethttp.Handler) gdt.Fixture {
	return &httpServerFixture{handler: h}
}
