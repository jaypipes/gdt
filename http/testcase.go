package http

import (
	"bytes"
	"io"
	"net/http"
	nethttp "net/http"
	"strings"
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type httpTestcaseConfig struct {
	baseURL string
}

// httpTestcase interfaces.Testcase
type httpTestcase struct {
	*gdt.TestFile
	cfg *httpTestcaseConfig
	// cache of last HTTP response one of the test units executed
	PrevResponse *nethttp.Response
}

func (htc *httpTestcase) BaseURL() string {
	if htc.cfg != nil && htc.cfg.baseURL != "" {
		return htc.cfg.baseURL
	}
	// query the fixture registry to determine if any of them contain an
	// http.baseurl state attribute.
	for _, f := range htc.Fixtures() {
		if f.HasState("http.base_url") {
			return f.State("http.base_url")
		}
	}
	return ""
}

// httpTest implements interfaces.Testable
type httpTest struct {
	tc *httpTestcase
	// Name for the individual HTTP call test
	name string
	// Description of the test (defaults to Name)
	description string
	// URL being called by HTTP client
	url string
	// HTTP Method specified by HTTP client
	method string
	// HTTP request to execute
	request *nethttp.Request
	// JSON-marshaled payload to send in request
	jsonBody []byte
	// HTTP Response object to assert on
	response *response
	// Specification for expected response
	responseAssertion *responseAssertion
}

// Run executes the test described by the HTTP test
func (ht *httpTest) Run(t *testing.T) {
	var err error
	baseURL := ht.tc.BaseURL()
	var body io.Reader
	if ht.jsonBody != nil {
		body = bytes.NewReader(ht.jsonBody)
	}
	var urlStr string
	if strings.ToUpper(ht.url) == "$LOCATION" {
		if ht.tc.PrevResponse == nil {
			panic("test unit referenced $LOCATION before executing an HTTP request")
		}
		url, err := ht.tc.PrevResponse.Location()
		if err != nil {
			panic(err)
		}
		urlStr = url.String()
	} else {
		urlStr = baseURL + ht.url
	}
	ht.request, err = http.NewRequest(ht.method, urlStr, body)
	require.Nil(t, err)
	c := nethttp.DefaultClient
	resp, err := c.Do(ht.request)
	require.Nil(t, err)
	require.NotNil(t, resp, "Expected nil net/http:Response but got nil")
	ht.response = &response{resp}
	t.Run(ht.name, func(t *testing.T) {
		if ht.responseAssertion != nil {
			rspec := ht.responseAssertion
			if rspec.JSON != nil {
				if rspec.JSON.Length != nil {
					ht.assertJSONLength(t, *(rspec.JSON.Length))
				}
			}

			if rspec.Status != nil {
				ht.assertStatusCode(t, *(rspec.Status))
			}

			if len(rspec.Strings) > 0 {
				for _, exp := range rspec.Strings {
					ht.assertStringIn(t, exp)
				}
			}
		}
	})
	ht.tc.PrevResponse = resp
}

func (ht *httpTest) assertJSONLength(t *testing.T, exp uint) {
	t.Run("check JSON length", func(t *testing.T) {
		got := ht.response.JSON()
		assert.Equal(t, uint(len(got)), exp, "Expected HTTP response to have JSON length of %d but got %d", exp, len(got))
	})
}

func (ht *httpTest) assertStatusCode(t *testing.T, exp int) {
	t.Run("check HTTP status code", func(t *testing.T) {
		got := ht.response.StatusCode
		assert.Equal(t, exp, got, "Expected HTTP response to have status code of %d but got %d", exp, got)
	})
}

func (ht *httpTest) assertStringIn(t *testing.T, exp string) {
	t.Run("check HTTP status code", func(t *testing.T) {
		got := ht.response.Text()
		assert.Contains(t, got, exp, "Expected HTTP response to contain %s", exp)
	})
}
