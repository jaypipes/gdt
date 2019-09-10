package http

import (
	"encoding/json"
	nethttp "net/http"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	msgHTTPStatus          = "Expected HTTP response to have status code of %d but got %d"
	msgJSONLength          = "Expected HTTP response to have JSON length of %d but got %d"
	msgJSONUnmarshalError  = "Failed to unmarshal JSON: %s"
	msgJSONPathError       = "Failed to find JSONPath %s: %s"
	msgJSONPathStringValue = "Expected string value %s but got JSONPath value %v which is not convertable to string"
	msgStringInBody        = "Expected HTTP response to contain %s"
	msgHeaderIn            = "Expected HTTP header %s to be in response"
	msgHeaderValue         = "Expected HTTP header with value %s to be in response"
)

func assertHTTPStatusEqual(t *testing.T, r *nethttp.Response, exp int) {
	t.Helper()
	got := r.StatusCode
	assert.Equal(t, exp, got, msgHTTPStatus, exp, got)
}

func assertStringInBody(t *testing.T, r *nethttp.Response, b []byte, exp string) {
	t.Helper()
	assert.Contains(t, string(b), exp, msgStringInBody, exp)
}

func assertHeader(t *testing.T, r *nethttp.Response, exp string) {
	t.Helper()
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

func assertJSON(t *testing.T, r *nethttp.Response, b []byte, jspec *jsonAssertion) {
	t.Helper()
	if jspec.Length != nil {
		// An error may have been returned as plain/text. In this case, we
		// don't want to check the length of the JSON-serialized body
		if strings.HasPrefix(r.Header.Get("content-type"), "application/json") {
			assertJSONLen(t, r, b, *(jspec.Length))
		}
	}
	if len(jspec.Paths) > 0 {
		assertJSONPaths(t, r, b, jspec.Paths)
	}
}

func assertJSONLen(t *testing.T, r *nethttp.Response, b []byte, exp uint) {
	t.Helper()
	assert.Equal(t, exp, uint(len(b)), msgJSONLength, exp, len(b))
}

func assertJSONPaths(t *testing.T, r *nethttp.Response, b []byte, paths map[string]string) {
	t.Helper()
	v := interface{}(nil)
	err := json.Unmarshal(b, &v)
	require.Nil(t, err, msgJSONUnmarshalError, err)
	for path, expVal := range paths {
		assertJSONPath(t, r, path, expVal, v)
	}
}

func assertJSONPath(t *testing.T, r *nethttp.Response, path string, exp string, v interface{}) {
	t.Helper()
	got, err := jsonpath.Get(path, v)
	require.Nil(t, err, msgJSONPathError, path, err)
	gotStr, ok := got.(string)
	assert.True(t, ok, msgJSONPathStringValue, exp, got)
	assert.Equal(t, exp, gotStr)
}
