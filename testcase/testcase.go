package testcase

import (
	"io"

	"github.com/jaypipes/gdt/interfaces"
)

// Testcase describes the tests in a single gdt test file. Implements
// interfaces.Testcase
type testcase struct {
	// typ is the type of test case
	typ string
	// filepath is the filepath to the test file
	filepath string
	// name for the overall test
	name string
	// description of the test (defaults to name)
	description string
	// tests that may be run for this test case
	tests []interfaces.Runnable
}

// AppendRunnable appends a Runnable thing to the test case's list of Runnable
// things
func (tc *testcase) AppendRunnable(r interfaces.Runnable) {
	tc.tests = append(tc.tests, r)
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

// SetDescription sets the test case's longer description
func (tc *testcase) SetDescription(description string) {
	tc.description = description
}

// New returns a new `Testcase` for an HTTP test case. The function
// accepts zero or more `WithOption` values that affect the returned test
// case.
//
// Usage:
//
//   tc := testcase.New(testcase.Withname("books_api"))
func New(opts ...WithOption) *testcase {
	useOpts := mergeOptions(opts...)
	t := &testcase{}
	if useOpts.Description != "" {
		t.description = useOpts.Description
	}
	if useOpts.Name != "" {
		t.name = useOpts.Name
	}
	return t
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

// Run executes the elements of the Testcase
func (tc *testcase) Run(ow, ew io.Writer) interfaces.RunResult {
	merged := &runResult{
		succeeded: true,
		skipped:   false,
		errors:    []error{},
	}
	for _, r := range tc.tests {
		res := r.Run(ow, ew)
		merged.succeeded = merged.succeeded && res.OK()
	}
	return merged
}
