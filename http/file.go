package http

import (
	"bytes"
	"io"
	"net/http"
	nethttp "net/http"
	"strings"
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/stretchr/testify/require"
)

type httpTestFileConfig struct {
	baseURL string
}

// httpTestFile wraps gdt.TestFile and adapts it for specialization of groups
// of tests of HTTP APIs
type httpTestFile struct {
	*gdt.TestFile
	cfg *httpTestFileConfig
	// cache of last HTTP response one of the test units executed
	PrevResponse *nethttp.Response
}

// BaseURL returns the base URL to use when constructing HTTP requests
func (htf *httpTestFile) BaseURL() string {
	// If the httpTestFile has been manually configured and the configuration
	// contains a base URL, use that. Otherwise, check to see if there is a
	// fixture in the registry that has an "http.base_url" state key and use
	// that if found.
	if htf.cfg != nil && htf.cfg.baseURL != "" {
		return htf.cfg.baseURL
	}
	// query the fixture registry to determine if any of them contain an
	// http.base_url state attribute.
	for _, f := range htf.Fixtures() {
		if f.HasState(FIXTURE_STATE_KEY_BASE_URL) {
			return f.State(FIXTURE_STATE_KEY_BASE_URL)
		}
	}
	return ""
}

// httpTest represents a single HTTP request and response pair along with
// expectations/assertions for the response components
type httpTest struct {
	tf *httpTestFile
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

func (ht *httpTest) getURL() string {
	if strings.ToUpper(ht.url) == "$LOCATION" {
		if ht.tf.PrevResponse == nil {
			panic("test unit referenced $LOCATION before executing an HTTP request")
		}
		url, err := ht.tf.PrevResponse.Location()
		if err != nil {
			panic(err)
		}
		return url.String()
	}
	baseURL := ht.tf.BaseURL()
	return baseURL + ht.url
}

// Run executes the test described by the HTTP test
func (ht *httpTest) Run(t *testing.T) {
	var body io.Reader
	if ht.jsonBody != nil {
		body = bytes.NewReader(ht.jsonBody)
	}
	req, err := http.NewRequest(ht.method, ht.getURL(), body)
	require.Nil(t, err)
	c := nethttp.DefaultClient
	resp, err := c.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp, "Expected nil net/http:Response but got nil")
	t.Run(ht.name, func(t *testing.T) {
		if ht.responseAssertion != nil {
			rspec := ht.responseAssertion

			if rspec.Status != nil {
				assertHTTPStatusEqual(t, resp, *(rspec.Status))
			}

			if rspec.JSON != nil {
				if rspec.JSON.Length != nil {
					assertJSONLen(t, resp, *(rspec.JSON.Length))
				}
			}

			if len(rspec.Strings) > 0 {
				for _, exp := range rspec.Strings {
					assertStringInBody(t, resp, exp)
				}
			}

			if len(rspec.Headers) > 0 {
				for _, exp := range rspec.Headers {
					assertHeader(t, resp, exp)
				}
			}
		}
	})
	ht.tf.PrevResponse = resp
}
