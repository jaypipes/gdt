// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"testing"
)

// suite contains zero or more Runnable things, one for each YAML file in a
// given directory
type suite struct {
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
func (s *suite) Append(r Runnable) {
	s.units = append(s.units, r)
}

// Run executes the tests in the test case
func (s *suite) Run(t *testing.T, ctx *Context) {
	for _, unit := range s.units {
		unit.Run(t, ctx)
	}
}
