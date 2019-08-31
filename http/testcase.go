package http

import (
	"net/http"
	nethttp "net/http"
	"strings"
	"testing"

	"github.com/jaypipes/gdt/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testUnit implements interfaces.Runnable
type testUnit struct {
	t *testing.T
	// fixture registry associated with the test unit's test case
	fr interfaces.FixtureRegistry
	// Name for the individual HTTP call test
	name string
	// Description of the test (defaults to Name)
	description string
	// Base URL to use for request
	baseURL string
	// URL being called by HTTP client
	URL string `json:"url"`
	// HTTP Method specified by HTTP client
	Method string `json:"method"`
	// HTTP request to execute
	request *nethttp.Request
	// HTTP Response object to assert on
	response *response
	// Specification for expected response
	responseAssertion *responseAssertion
}

// T returns the underlying pointer to a testing.T
func (tu *testUnit) T() *testing.T {
	return tu.t
}

// Name returns a name for the test unit
func (tu *testUnit) Name() string {
	return tu.name
}

// Describe returns a description or name for the test unit
func (tu *testUnit) Describe() string {
	return tu.description
}

// Run executes the test described by the test unit
func (tu *testUnit) RunWithFixtures(fr interfaces.FixtureRegistry) {
	var err error
	baseURL := fr.Get("books_api").State("URL")
	tu.request, err = http.NewRequest(tu.Method, baseURL+tu.URL, nil)
	require.Nil(tu.t, err)
	c := nethttp.DefaultClient
	resp, _ := c.Do(tu.request)
	tu.response = &response{resp}
	require.NotNil(tu.t, resp, "Expected nil net/http:Response but got nil")
	tu.t.Run(tu.name, func(t *testing.T) {
		if tu.responseAssertion != nil {
			rspec := tu.responseAssertion
			if rspec.JSON != nil {
				if rspec.JSON.Length != nil {
					tu.assertJSONLength(*(rspec.JSON.Length))
				}
			}

			if rspec.Status != nil {
				tu.assertStatusCode(*(rspec.Status))
			}

			if len(rspec.Strings) > 0 {
				for _, exp := range rspec.Strings {
					tu.assertStringIn(exp)
				}
			}
		}
	})
}

func (tu *testUnit) requestURL(path string) string {
	return tu.baseURL + "/" + strings.TrimPrefix(path, "/")
}

func (tu *testUnit) assertJSONLength(exp uint) {
	tu.t.Run("check JSON length", func(t *testing.T) {
		got := tu.response.JSON()
		assert.Equal(t, uint(len(got)), exp, "Expected HTTP response to have JSON length of %d but got %d", exp, len(got))
	})
}

func (tu *testUnit) assertStatusCode(exp int) {
	tu.t.Run("check HTTP status code", func(t *testing.T) {
		got := tu.response.StatusCode
		assert.Equal(t, exp, got, "Expected HTTP response to have status code of %d but got %d", exp, got)
	})
}

func (tu *testUnit) assertStringIn(exp string) {
	tu.t.Run("check HTTP status code", func(t *testing.T) {
		got := tu.response.Text()
		assert.Contains(t, got, exp, "Expected HTTP response to contain %s", exp)
	})
}
