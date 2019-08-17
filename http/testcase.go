package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jaypipes/gdt"
	"github.com/jaypipes/gdt/testcase"
)

// HTTPTest implements gdt.Runnable
type HTTPTest struct {
	// Name for the individual HTTP call test
	name string
	// Description of the test (defaults to Name)
	description string
	// Base URL to use for request
	baseURL string
	// HTTP request to execute
	request *http.Request
	// HTTP Response object to assert on
	response *response
	// Specification for expected response
	assertions []gdt.Assertion
}

func (t *HTTPTest) requestURL() string {
	return t.baseURL + "/" + strings.TrimPrefix(t.path, "/")
}

func assertZeroJSONLength(resp *http.Response) func(*http.Response) (bool, string) {
	return func(resp *http.Response) (bool, string) {
		res := respJSON(resp)
		if res != "" {
			return false, fmt.Sprintf("Expected HTTP response to have no JSON but found %s", res)
		}
		return true, ""
	}
}

func (t *HTTPTest) Run() gdt.RunResult {
	succeeded := true
	skipped := false
	errs := make([]error, 0)
	t.response, err = http.Get(apiPath(t.url))

	for a := range t.assertions {
		succeeded &= a.Assert()
	}
	return gdt.RunResult{
		Succeeded: succeeded,
		Skipped:   skipped,
		Errors:    errs,
	}
}

func NewFromYAML(contents string, opts ...gdt.WithOption) *gdt.TestCase {
	tcs, err := parseYAML(contents)
	if err != nil {
		return nil, err
	}

	if tcs.Name != "" {
		opts = append(opts, gdt.WithName(tcs.Name))
	}
	if tcs.Description != "" {
		opts = append(opts, gdt.WithDescription(tcs.Description))
	}

	tc := testcase.New(opts...)
	for _, tspec := range tcs.TestSpecs {
		ht := HTTPTest{
			name: tspec.Name,
		}

		if tspec.URL == "" {
			if tspec.GET != "" {
				ht.request = http.Request(URL: tspec.GET, Method: "GET")
			}
		
		} else {
			method := tspec.Method
			if method == "" {
				return nil, fmt.Errorf("When specifying url in HTTP test spec, please specify an HTTP method")
			}

		}

		if tspec.Response != nil {
			rspec := tspec.Response
			if rspec.JSON != nil {
				if rspec.JSON.Length == 0 {
					r.AssertZeroJSONLength()
				}
			}

			if rspec.Status != 0 {
				It(fmt.Sprintf("should return %d", rspec.Status), func() {
					Ω(response.StatusCode).Should(Equal(rspec.Status))
				})
			}

			if len(rspec.Strings) > 0 {
				for _, expStr := range rspec.Strings {
					It(fmt.Sprintf("should contain '%s'", expStr), func() {
						Ω(respText(response)).Should(ContainSubstring(expStr))
					})
				}
			}
		}
	}
}
