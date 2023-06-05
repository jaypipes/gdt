// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"testing"
)

// TestSuite contains zero or more Runnable things, one for each YAML file in a
// given directory
type TestSuite struct {
	// Name is the short name for the test suite. If empty, defaults to Path
	Name string `json:"name,omitempty"`
	// Description is a description of the tests contained in the test suite
	Description string `json:"description,omitempty"`
	// path is the filepath to the test suite directory
	path string
	// Collection of test files in this test suite
	units []Runnable
}

// Append appends a runnable test element to the test suite
func (s *TestSuite) Append(r Runnable) {
	s.units = append(s.units, r)
}

// Run executes the tests in the test case
func (s *TestSuite) Run(t *testing.T, ctx *Context) {
	for _, unit := range s.units {
		unit.Run(t, ctx)
	}
}

// NewTestSuite returns a new TestSuite
func NewTestSuite(options ...*Option) *TestSuite {
	merged := mergeOptions(options)
	s := &TestSuite{}
	if merged.path != nil {
		s.path = *merged.path
	}
	if merged.name != nil {
		s.Name = *merged.name
	}
	if merged.description != nil {
		s.Description = *merged.description
	}
	return s
}
