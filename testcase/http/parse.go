package http

import (
	"fmt"
	"net/http"

	"github.com/ghodss/yaml"
)

type jsonAssertion struct {
	Length *uint `json:"length"`
}

type responseAssertion struct {
	JSON    *jsonAssertion `json:"json"`
	Strings []string       `json:"strings"`
	Status  *int           `json:"status"`
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
	GET string `json:"GET"`
	// Shortcut for URL and Method of "POST"
	POST string `json:"POST"`
	// Specification for expected response
	Response *responseAssertion `json:"response"`
}

type testcaseSpec struct {
	Specs []*testSpec `json:"tests"`
}

// Parse accepts a Testcase and a string of YAML contents from a gdt test file.
// It then parses the HTTP test case and adds the HTTP-specific tests to the
// supplied Testcase
func Parse(tc *test.Testcase, contents string) error {
	tcs := testcaseSpec{}
	if err := yaml.Unmarshal(contents, &tcs); err != nil {
		return err
	}
	for _, tspec := range tcs.TestSpecs {
		ht := HTTPTest{
			name: tspec.Name,
		}

		if tspec.URL == "" {
			if tspec.GET != "" {
				ht.request = http.NewRequest(tspec.GET, "GET")
			}
		} else {
			method := tspec.Method
			if method == "" {
				return nil, fmt.Errorf("When specifying url in HTTP test spec, please specify an HTTP method")
			}
			ht.request = http.NewRequest(tspec.URL, method)
		}

		if tspec.Response != nil {
			rspec := tspec.Response
			if rspec.JSON != nil {
				if rspec.JSON.Length != nil {
					ht.assertJSONLength(*rspec.JSON.Length)
				}
			}

			if rspec.Status != nil {
				ht.assertStatusCode(*rspec.Status)
			}

			if len(rspec.Strings) > 0 {
				for _, exp := range rspec.Strings {
					ht.assertStringIn(exp)
				}
			}
		}
		tc.AppendRunnable(&ht)
	}
	return nil
}
