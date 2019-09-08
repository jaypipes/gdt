package http

import (
	nethttp "net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	msgHTTPStatus   = "Expected HTTP response to have status code of %d but got %d"
	msgJSONLength   = "Expected HTTP response to have JSON length of %d but got %d"
	msgStringInBody = "Expected HTTP response to contain %s"
	msgHeaderIn     = "Expected header %s to be in response"
	msgHeaderValue  = "Expected header with value %s to be in response"
)

func assertHTTPStatusEqual(t *testing.T, r *nethttp.Response, exp int) {
	t.Helper()
	got := r.StatusCode
	assert.Equal(t, exp, got, msgHTTPStatus, exp, got)
}

func assertJSONLen(t *testing.T, r *nethttp.Response, exp uint) {
	t.Helper()
	response := response{r}
	got := response.JSON()
	assert.Equal(t, uint(len(got)), exp, msgJSONLength, exp, len(got))
}

func assertStringInBody(t *testing.T, r *nethttp.Response, exp string) {
	response := response{r}
	got := response.Text()
	assert.Contains(t, got, exp, msgStringInBody, exp)
}

func assertHeader(t *testing.T, r *nethttp.Response, exp string) {
	colonPos := strings.IndexRune(exp, ':')
	if colonPos > -1 {
		keyPart := exp[:colonPos]
		valPart := exp[colonPos+1:]
		val := r.Header.Get(keyPart)
		assert.NotEmpty(t, val, msgHeaderIn, exp)
		// If the string being compared is of the form Key: Value,
		// then we check for both existence and the value of the
		// header
		expVal := strings.ToLower(valPart)
		assert.Equal(t, expVal, strings.ToLower(val), msgHeaderValue, exp)
	} else {
		val := r.Header.Get(exp)
		assert.NotEmpty(t, val, msgHeaderIn, exp)
	}
}
