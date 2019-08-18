package suite

import "../testcase"

// Suite contains TestCases
type Suite struct {
	// Filepath is the filepath to the test suite file
	Filepath string
	// Name for the test suite (defaults to Filepath
	Name string
	// Description of the test suite (defaults to Name)
	Description string
	// Collection of test cases in this suite
	TestCases []*testcase.TestCase
}
