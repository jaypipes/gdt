package testcase

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type RunResult struct {
	Succeeded bool
	Skipped bool
	Errors []error
}

type Runnable interface {
	func Run() RunResult
}

type Assertion interface {
	func Assert(w io.Writer) bool
}

type TestCaseType uint

const (
	TestCaseTypeHTTP TestCaseType = iota
)

type TestCase struct {
	// Type of test case this test has
	TestCaseType TestCaseType `json:"type"`
	// Filepath is the filepath to the test file
	Filepath string `json:"filepath"`
	// Name for the overall test
	Name string `json:"name"`
	// Description of the test (defaults to Name)
	Description string `json:"description"`
	// Tests that may be run for this test case
	Tests []Runnable string `json:"tests"`
}

func (tc *TestCase) AppendRunnable(r Runnable) {
	tc.Tests = append(tc.Tests, r)
}

// TestSuite contains TestCases
type TestSuite struct {
	TestCases []*TestCase `json:"test_cases"`
}
