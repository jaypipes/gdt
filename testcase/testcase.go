package testcase

import (
	"testing"
)

// Testcase describes the tests in a single gdt test file. Implements
// interfaces.Testcase and wraps the testing.T struct
type testcase struct {
	t *testing.T
	// typ is the type of test case
	typ string
	// filepath is the filepath to the test file
	filepath string
	// name for the overall test
	name string
	// description of the test (defaults to name)
	description string
	// set of fixture names and args to associate with the testcase's
	// before-run stage
	before map[string][]string
}

// T returns a pointer to the testing.T
func (tc *testcase) T() *testing.T {
	return tc.t
}

// Type returns the test case's type, e.g. "http"
func (tc *testcase) Type() string {
	return tc.typ
}

// Filepath returns the test case's absolute filepath
func (tc *testcase) Filepath() string {
	return tc.filepath
}

// Name returns a name for the test case
func (tc *testcase) Name() string {
	return tc.name
}

// Describe returns a description or name for the test case
func (tc *testcase) Describe() string {
	return tc.description
}

// New returns a new `Testcase` for an HTTP test case. The function
// accepts zero or more `WithOption` values that affect the returned test
// case.
//
// Usage:
//
//   tc := testcase.New(testcase.Withname("books_api"))
func New(t *testing.T, opts ...WithOption) *testcase {
	useOpts := mergeOptions(opts...)
	tc := &testcase{t: t}
	if useOpts.Description != "" {
		tc.description = useOpts.Description
	}
	if useOpts.Name != "" {
		tc.name = useOpts.Name
	}
	return tc
}
