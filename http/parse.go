package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/jaypipes/gdt/interfaces"
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

type httpParser struct{}

// Parse accepts a Testcase and a string of YAML contents from a gdt test file.
// It then parses the HTTP test case and adds the HTTP-specific tests to the
// supplied Testcase
func (p *httpParser) Parse(tc interfaces.Testcase, contents []byte) error {
	var err error
	tcs := testcaseSpec{}
	if err := yaml.Unmarshal(contents, &tcs); err != nil {
		return err
	}
	for _, tspec := range tcs.Specs {
		tu := testUnit{
			name: tspec.Name,
		}

		if tspec.URL == "" {
			if tspec.GET != "" {
				tu.request, err = http.NewRequest("GET", tspec.GET, strings.NewReader(""))
				if err != nil {
					return err
				}
			}
		} else {
			method := tspec.Method
			if method == "" {
				return fmt.Errorf("When specifying url in HTTP test spec, please specify an HTTP method")
			}
			tu.request, err = http.NewRequest(method, tspec.URL, strings.NewReader(""))
			if err != nil {
				return err
			}
		}

		if tspec.Response != nil {
			rspec := tspec.Response
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
		tc.AppendRunnable(&tu)
	}
	return nil
}
