package testcase

type RunResult struct {
	Succeeded bool
	Skipped   bool
	Errors    []error
}

type TestCaseType uint

const (
	TestCaseTypeHTTP TestCaseType = iota
)

type TestCase struct {
	// Type of test case this test has
	TestCaseType TestCaseType
	// Filepath is the filepath to the test file
	Filepath string
	// Name for the overall test
	Name string
	// Description of the test (defaults to Name)
	Description string
	// Tests that may be run for this test case
	Tests []gdt.Runnable
	// Fixtures provide extensible setup and teardown of resources used in the
	// test case's tests
	Fixtures []gdt.Fixture
}

func (tc *TestCase) AppendRunnable(r Runnable) {
	tc.Tests = append(tc.Tests, r)
}

// TestSuite contains TestCases
type TestSuite struct {
	TestCases []*TestCase `json:"test_cases"`
}
