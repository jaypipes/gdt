package http

import (
	"net/http"
	nethttp "net/http"
	"testing"

	"github.com/jaypipes/gdt/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type httpTestcaseConfig struct {
	baseURL string
}

// httpTestcase interfaces.Testcase
type httpTestcase struct {
	interfaces.Testcase
	cfg *httpTestcaseConfig
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
	url string `json:"url"`
	// HTTP Method specified by HTTP client
	method string `json:"method"`
	// HTTP request to execute
	request *nethttp.Request
	// HTTP Response object to assert on
	response *response
	// Specification for expected response
	responseAssertion *responseAssertion
}

// Name returns a name for the HTTP test
func (ht *httpTest) Name() string {
	return ht.name
}

// Describe returns a description or name for the HTTP test
func (ht *httpTest) Describe() string {
	return ht.description
}

// Run executes the test described by the HTTP test
func (ht *httpTest) Run() {
	var err error
	baseURL := ht.tc.BaseURL()
	ht.request, err = http.NewRequest(ht.method, baseURL+ht.url, nil)
	require.Nil(ht.tc.T(), err)
	c := nethttp.DefaultClient
	resp, _ := c.Do(ht.request)
	ht.response = &response{resp}
	require.NotNil(ht.tc.T(), resp, "Expected nil net/http:Response but got nil")
	ht.tc.T().Run(ht.name, func(t *testing.T) {
		if ht.responseAssertion != nil {
			rspec := ht.responseAssertion
			if rspec.JSON != nil {
				if rspec.JSON.Length != nil {
					ht.assertJSONLength(*(rspec.JSON.Length))
				}
			}

			if rspec.Status != nil {
				ht.assertStatusCode(*(rspec.Status))
			}

			if len(rspec.Strings) > 0 {
				for _, exp := range rspec.Strings {
					ht.assertStringIn(exp)
				}
			}
		}
	})
}

func (tu *httpTest) assertJSONLength(exp uint) {
	tu.tc.T().Run("check JSON length", func(t *testing.T) {
		got := tu.response.JSON()
		assert.Equal(t, uint(len(got)), exp, "Expected HTTP response to have JSON length of %d but got %d", exp, len(got))
	})
}

func (tu *httpTest) assertStatusCode(exp int) {
	tu.tc.T().Run("check HTTP status code", func(t *testing.T) {
		got := tu.response.StatusCode
		assert.Equal(t, exp, got, "Expected HTTP response to have status code of %d but got %d", exp, got)
	})
}

func (tu *httpTest) assertStringIn(exp string) {
	tu.tc.T().Run("check HTTP status code", func(t *testing.T) {
		got := tu.response.Text()
		assert.Contains(t, got, exp, "Expected HTTP response to contain %s", exp)
	})
}
