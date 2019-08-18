package testcase

import (
	"errors"
	"io"
)

// Type of test case
type Type uint

const (
	// TypeHTTP describes tests of HTTP APIs
	TypeHTTP Type = iota
)

var (
	// ErrNilTestcase is returned when attempting to call a Testcase method
	// with a nil Testcase pointer
	ErrNilTestcase = errors.New("nil testcase")
)

// RunResult is what is returned from a Runnable
type RunResult struct {
	Succeeded bool
	Skipped   bool
	Errors    []error
}

// Runnable describes the interface that a Testcase runs
type Runnable interface {
	// Run takes two io.Writers, one for the normal output stream, the other
	// for the error stream, and returns a RunResult
	Run(ow, ew io.Writer) RunResult
}

// Testcase describes the tests in a single gdt test file
type Testcase struct {
	// Type of test case this test has
	Type Type
	// Filepath is the filepath to the test file
	Filepath string
	// Name for the overall test
	Name string
	// Description of the test (defaults to Name)
	Description string
	// Tests that may be run for this test case
	Tests []Runnable
	// Fixtures provide extensible setup and teardown of resources used in the
	// test case's tests
	Fixtures []struct{}
}

// AppendRunnable appends a Runnable thing to the test case's list of Runnable
// things
func (tc *Testcase) AppendRunnable(r Runnable) {
	tc.Tests = append(tc.Tests, r)
}

// New returns a new `Testcase` for an HTTP test case. The function
// accepts zero or more `WithOption` values that affect the returned test
// case.
//
// Usage:
//
//   tc := testcase.New(testcase.WithName("books_api"))
func New(opts ...WithOption) *Testcase {
	useOpts := mergeOptions(opts...)
	t := &Testcase{
		Type: TypeHTTP,
	}
	if useOpts.Description != "" {
		t.Description = useOpts.Description
	}
	if useOpts.Filepath != "" {
		t.Filepath = useOpts.Filepath
	}
	if useOpts.Name != "" {
		t.Name = useOpts.Name
	}
	return t
}

// Run executes the elements of the Testcase
func (tc *Testcase) Run(ow, ew io.Writer) []RunResult {
	results := make([]RunResult, len(tc.Tests))
	for x, r := range tc.Tests {
		res[x] = r.Run(ow, ew)
	}
	return results
}
