package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	nethttp "net/http"
	"strings"
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/stretchr/testify/require"
)

type httpFileConfig struct {
	baseURL string
}

// httpFile contains groups of tests of HTTP APIs
type httpFile struct {
	ctx *gdt.Context
	cfg *httpFileConfig
	// cache of last HTTP response one of the test units executed
	PrevResponse *nethttp.Response
}

// BaseURL returns the base URL to use when constructing HTTP requests
func (hf *httpFile) BaseURL() string {
	// If the httpFile has been manually configured and the configuration
	// contains a base URL, use that. Otherwise, check to see if there is a
	// fixture in the registry that has an "http.base_url" state key and use
	// that if found.
	if hf.cfg != nil && hf.cfg.baseURL != "" {
		return hf.cfg.baseURL
	}
	// query the fixture registry to determine if any of them contain an
	// http.base_url state attribute.
	for _, f := range hf.ctx.Fixtures.List() {
		if f.HasState(FIXTURE_STATE_KEY_BASE_URL) {
			return f.State(FIXTURE_STATE_KEY_BASE_URL)
		}
	}
	return ""
}

// httpTest represents a single HTTP request and response pair along with
// expectations/assertions for the response components
type httpTest struct {
	f *httpFile
	// Name for the individual HTTP call test
	name string
	// Description of the test (defaults to Name)
	description string
	// URL being called by HTTP client
	url string
	// HTTP Method specified by HTTP client
	method string
	// JSON-marshaled payload to send in request
	jsonBody []byte
	// Specification for expected response
	responseAssertion *responseAssertion
}

// getURL returns the URL to use for the test's HTTP request. The test's url
// field is first queried to see if it is the special $LOCATION string. If it
// is, then we return the previous HTTP response's Location header. Otherwise,
// we construct the URL from the httpFile's base URL and the test's url field.
func (ht *httpTest) getURL() string {
	if strings.ToUpper(ht.url) == "$LOCATION" {
		if ht.f.PrevResponse == nil {
			panic("test unit referenced $LOCATION before executing an HTTP request")
		}
		url, err := ht.f.PrevResponse.Location()
		if err != nil {
			panic(err)
		}
		return url.String()
	}
	baseURL := ht.f.BaseURL()
	return baseURL + ht.url
}

// Run executes the test described by the HTTP test. A new HTTP request and
// response pair is created during this call.
func (ht *httpTest) Run(t *testing.T) {
	var body io.Reader
	if ht.jsonBody != nil {
		body = bytes.NewReader(ht.jsonBody)
	}
	t.Run(ht.name, func(t *testing.T) {
		req, err := http.NewRequest(ht.method, ht.getURL(), body)
		require.Nil(t, err)
		// TODO(jaypipes): Allow customization of the HTTP client for proxying,
		// TLS, etc
		c := nethttp.DefaultClient
		resp, err := c.Do(req)
		require.Nil(t, err)
		if ht.responseAssertion != nil {
			// Only read the response body contents once and pass the byte
			// buffer to the assertion functions
			b, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)

			rspec := ht.responseAssertion
			if rspec.Status != nil {
				assertHTTPStatusEqual(t, resp, *(rspec.Status))
			}

			if rspec.JSON != nil {
				assertJSON(t, resp, b, rspec.JSON)
			}

			if len(rspec.Strings) > 0 {
				for _, exp := range rspec.Strings {
					assertStringInBody(t, resp, b, exp)
				}
			}

			if len(rspec.Headers) > 0 {
				for _, exp := range rspec.Headers {
					assertHeader(t, resp, exp)
				}
			}
		}
		ht.f.PrevResponse = resp
	})
}
