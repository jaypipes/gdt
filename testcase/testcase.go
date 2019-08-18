package testcase

// Type of test case
type Type uint

const (
	// TypeHTTP describes tests of HTTP APIs
	TypeHTTP Type = iota
)

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
	Tests []struct{}
	// Fixtures provide extensible setup and teardown of resources used in the
	// test case's tests
	Fixtures []struct{}
}
