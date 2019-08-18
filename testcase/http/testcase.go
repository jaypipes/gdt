package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"../../testcase"
)

// HTTPTest implements testcase.Runnable
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
	assertions []testcase.Assertion
}

func (t *HTTPTest) requestURL(path string) string {
	return t.baseURL + "/" + strings.TrimPrefix(path, "/")
}

type httpAssertion struct {
	test       *HTTPTest
	comparator func(r *response) (bool, string)
}

func (ha httpAssertion) Assert() (bool, string) {
	return comparator(t.response)
}

func (t *HTTPTest) assertJSONLength(exp int) {
	t.assertions = append(t.assertions, httpAssertion{
		test: t,
		comparator: func(r *response) (bool, string) {
			got := r.JSON()
			if len(got) != exp {
				return false, fmt.Sprintf("Expected HTTP response to have JSON length of %d but got %d", exp, len(got))
			}
			return true, ""
		},
	})
}

func (t *HTTPTest) assertStatusCode(exp int) {
	t.assertions = append(t.assertions, httpAssertion{
		test: t,
		comparator: func(r *response) (bool, string) {
			got := r.StatusCode
			if got != exp {
				return false, fmt.Sprintf("Expected HTTP response to have status code of %d but got %d", exp, got)
			}
			return true, ""
		},
	})
}

func (t *HTTPTest) assertStringIn(exp string) {
	t.assertions = append(t.assertions, httpAssertion{
		test: t,
		comparator: func(r *response) (bool, string) {
			got := r.Text()
			if strings.Contains(got, exp) {
				return false, fmt.Sprintf("Expected HTTP response to contain %s", exp)
			}
			return true, ""
		},
	})
}

func (t *HTTPTest) Run(_, _ io.Writer) testcase.RunResult {
	succeeded := true
	skipped := false
	errs := make([]error, 0)
	c := http.DefaultClient
	resp, err := c.Do(t.request)
	if err != nil {
		errs = append(errs, err)
	} else {
		t.response = response{resp}
		for a := range t.assertions {
			ok, failStr := a.Assert()
			succeeded &= ok
			if !ok && failStr != "" {
				errs = append(errs, errors.New(failStr))
			}
		}
	}
	return testcase.RunResult{
		Succeeded: succeeded,
		Skipped:   skipped,
		Errors:    errs,
	}
}
