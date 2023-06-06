// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package http

import (
	"encoding/json"
	nethttp "net/http"
	"strings"
	"testing"

	"github.com/PaesslerAG/jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gjs "github.com/xeipuuv/gojsonschema"
)

const (
	msgHTTPStatus          = "Expected HTTP response to have status code of %d but got %d"
	msgJSONLength          = "Expected HTTP response to have JSON length of %d but got %d"
	msgJSONUnmarshalError  = "Failed to unmarshal JSON: %s"
	msgJSONPathError       = "Failed to find JSONPath %s: %s"
	msgJSONPathStringValue = "Expected string value at JSONPath %s but got JSONPath value %v which is not convertable to string"
	msgJSONSchemaInvalid   = "Expected JSON to validate against JSONSchema, but found validation errors:\n%s"
	msgStringInBody        = "Expected HTTP response to contain %s"
	msgHeaderIn            = "Expected HTTP header %s to be in response"
	msgHeaderValue         = "Expected HTTP header with value %s to be in response"
	msgFormatInvalid       = "Unknown format %s in test"
	msgFormatBad           = "Expected %s to be formatted as %s"
)

// JSONAssertion represents one or more assertions about JSON data responses
type JSONAssertion struct {
	Length      *uint             `json:"length,omitempty"`
	Paths       map[string]string `json:"paths,omitempty"`
	PathFormats map[string]string `json:"path_formats,omitempty"`
	Schema      string            `json:"schema,omitempty"`
}

// ResponseAssertion contains one or more assertions about an HTTP response
type ResponseAssertion struct {
	// JSON contains the assertions about JSON data in the response
	JSON *JSONAssertion `json:"json,omitempty"`
	// Headers contains a list of HTTP headers that should be in the response
	Headers []string `json:"headers,omitempty"`
	// Strings contains a list of strings that should be present in the
	// response content
	Strings []string `json:"strings,omitempty"`
	// Status contains the numeric HTTP status code (e.g. 200 or 404) that
	// should be returned in the HTTP response
	Status *int `json:"status,omitempty"`
}

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

func assertJSON(t *testing.T, r *nethttp.Response, b []byte, jspec *JSONAssertion) {
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
	if len(jspec.PathFormats) > 0 {
		assertJSONPathFormats(t, r, b, jspec.PathFormats)
	}
	if jspec.Schema != "" {
		assertJSONSchema(t, r, b, jspec.Schema)
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
	assert.True(t, ok, msgJSONPathStringValue, path, got)
	assert.Equal(t, exp, gotStr)
}

func assertJSONPathFormats(t *testing.T, r *nethttp.Response, b []byte, pathFormats map[string]string) {
	t.Helper()
	v := interface{}(nil)
	err := json.Unmarshal(b, &v)
	require.Nil(t, err, msgJSONUnmarshalError, err)
	for path, format := range pathFormats {
		assertJSONPathFormat(t, r, path, format, v)
	}
}

func assertJSONPathFormat(t *testing.T, r *nethttp.Response, path string, format string, v interface{}) {
	t.Helper()
	got, err := jsonpath.Get(path, v)
	require.Nil(t, err, msgJSONPathError, path, err)
	gotStr, ok := got.(string)
	assert.True(t, ok, msgJSONPathStringValue, path, got)
	ok, err = isFormatted(format, got)
	require.Nil(t, err, msgFormatInvalid, format)
	assert.True(t, ok, msgFormatBad, format, gotStr)
}

// assertJSONSchema verifies that the HTTP response validates against a
// supplied JSONSchema document
//
// NOTE(jaypipes): schemaPath is an absolute path and should be checked for
// existence before running this function
func assertJSONSchema(
	t *testing.T,
	r *nethttp.Response,
	subject []byte,
	schemaPath string,
) {
	t.Helper()

	schemaLoader := gjs.NewReferenceLoader(schemaPath)
	docLoader := gjs.NewStringLoader(string(subject))

	res, err := gjs.Validate(schemaLoader, docLoader)
	require.Nil(t, err)

	var errStr string
	if len(res.Errors()) > 0 {
		errStrs := make([]string, len(res.Errors()))
		for x, e := range res.Errors() {
			errStrs[x] = e.String()
		}
		errStr = "- " + strings.Join(errStrs, "\n- ")
	}
	assert.True(t, res.Valid(), msgJSONSchemaInvalid, errStr)
}
