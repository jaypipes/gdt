package http

import (
	nethttp "net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testUnit implements interfaces.Runnable
type testUnit struct {
	t *testing.T
	// Name for the individual HTTP call test
	name string
	// Description of the test (defaults to Name)
	description string
	// Base URL to use for request
	baseURL string
	// HTTP request to execute
	request *nethttp.Request
	// HTTP Response object to assert on
	response *response
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
