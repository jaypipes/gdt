package http

import (
	"github.com/ghodss/yaml"
)

type JSONResponseAssertion struct {
	Length uint `json:"length"`
}

type HTTPResponseAssertion struct {
	JSON    *JSONResponseAssertion `json:"json"`
	Strings []string               `json:"strings"`
	Status  int                    `json:"status"`
}

type testSpec struct {
	// Name for the individual HTTP call test
	Name string `json:"name"`
	// Description of the test (defaults to Name)
	Description string `json:"description"`
	// URL being called by HTTP client
	URL string `json:"url"`
	// HTTP Method specified by HTTP client
	Method string `json:"method"`
	// Shortcut for URL and Method of "GET"
	GET string
	// Shortcut for URL and Method of "POST"
	POST string
	// HTTP request object constructed from spec
	Request *http.Request
	// Specification for expected response
	Response *HTTPResponseAssertion `json:"response"`
}

type testCaseSpec
	Specs []*testSpec `json:"tests"`
}

func parseYAML(contents string) (*testCaseSpec, error) {
	tcs := testCaseSpec{}
	if err := yaml.Unmarshal(contents, &tcs); err != nil {
		return nil, err
	}
	return &tcs, nil
}
