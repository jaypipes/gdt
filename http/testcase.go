package http

import (
	"errors"
	"fmt"
	"io"
	nethttp "net/http"
	"strings"

	"github.com/jaypipes/gdt/interfaces"
)

// testUnit implements interfaces.Runnable
type testUnit struct {
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
	// Specification for expected response
	assertions []interfaces.Assertion
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

type httpAssertion struct {
	test       *testUnit
	comparator func(r *response) (bool, string)
}

func (ha httpAssertion) Assert() (bool, string) {
	return ha.comparator(ha.test.response)
}

func (tu *testUnit) assertJSONLength(exp uint) {
	tu.assertions = append(tu.assertions, httpAssertion{
		test: tu,
		comparator: func(r *response) (bool, string) {
			got := r.JSON()
			if uint(len(got)) != exp {
				return false, fmt.Sprintf("Expected HTTP response to have JSON length of %d but got %d", exp, len(got))
			}
			return true, ""
		},
	})
}

func (tu *testUnit) assertStatusCode(exp int) {
	tu.assertions = append(tu.assertions, httpAssertion{
		test: tu,
		comparator: func(r *response) (bool, string) {
			got := r.StatusCode
			if got != exp {
				return false, fmt.Sprintf("Expected HTTP response to have status code of %d but got %d", exp, got)
			}
			return true, ""
		},
	})
}

func (tu *testUnit) assertStringIn(exp string) {
	tu.assertions = append(tu.assertions, httpAssertion{
		test: tu,
		comparator: func(r *response) (bool, string) {
			got := r.Text()
			if strings.Contains(got, exp) {
				return false, fmt.Sprintf("Expected HTTP response to contain %s", exp)
			}
			return true, ""
		},
	})
}

// runResult implements interfaces.RunResult
type runResult struct {
	succeeded bool
	skipped   bool
	errors    []error
}

func (r *runResult) OK() bool {
	return r.succeeded
}

func (r *runResult) Skipped() bool {
	return r.skipped
}

func (r *runResult) Errors() []error {
	return r.errors
}

// Run executes the HTTP call and returns the results of the client call
func (t *testUnit) Run(_, _ io.Writer) interfaces.RunResult {
	succeeded := true
	skipped := false
	errs := make([]error, 0)
	c := nethttp.DefaultClient
	resp, err := c.Do(t.request)
	if err != nil {
		errs = append(errs, err)
	} else {
		t.response = &response{resp}
		for _, a := range t.assertions {
			ok, failStr := a.Assert()
			succeeded = succeeded && ok
			if !ok && failStr != "" {
				errs = append(errs, errors.New(failStr))
			}
		}
	}
	return &runResult{
		succeeded: succeeded,
		skipped:   skipped,
		errors:    errs,
	}
}
