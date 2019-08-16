package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/jaypipes/gdt"
	"github.com/jaypipes/gdt/testcase"
)

type httpResponseAssertion struct {
	check func(*http.Response) (bool, string)
}

func (a *httpResponseAssertion) Assert(w io.Writer) bool {
	res, failMsg := a.check(a.resp)
	if !res {
		w.Write(failMsg)
	}
	return res
}

// HTTPTest implements gdt.Runnable
type HTTPTest struct {
	// Name for the individual HTTP call test
	name string
	// Description of the test (defaults to Name)
	description string
	// URL to query
	url string
	// HTTP Response object to assert on
	response *http.Response
	// Specification for expected response
	assertions []gdt.Assertion
}

func (t *HTTPTest) AssertZeroJSONLength() {
	t.assertions = append(t.assertions, httpResponseAssertion{assertZeroJSONLength(t.response)})
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

func NewTestCaseFromYAML(contents string, opts ...gdt.WithOption) *gdt.TestCase {
	tcs, err := parseYAML(contents)
	if err != nil {
		return nil, err
	}

	if tcs.Name != "" {
		opts = append(opts, gdt.WithName(tcs.Name))
	}
	if tcs.Description != "" {
		opts = append(opts, gdt.WithName(tcs.Description))
	}

	tc := testcase.New(opts...)
	for _, tspec := range ht.TestSpecs {
		r := HTTPTest{
			name: tspec.Name,
			url:  tspec.GET,
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
