package gdt

import (
	"testing"
)

// TestSuite contains zero or more Runnable things, one for each YAML file in a
// given directory
type TestSuite struct {
	// path is the filepath to the test suite directory
	path string
	// Name for the test suite
	name string
	// Description of the test suite
	description string
	// Collection of test files in this test suite
	units []Runnable
}

// Append appends a runnable test element to the test suite
func (ts *TestSuite) Append(r Runnable) {
	ts.units = append(ts.units, r)
}

// Run executes the tests in the test case
func (ts *TestSuite) Run(t *testing.T) {
	for _, unit := range ts.units {
		unit.Run(t)
	}
}
