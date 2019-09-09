package gdt

import (
	"testing"
)

// file describes the tests in a single gdt test file. Implements
// interfaces.Appendable and interfaces.Runnable
type file struct {
	ctx *Context
	// typ is the type of test in the file
	typ string
	// path is the filepath to the file
	path string
	// name for the overall test
	name string
	// description of the test (defaults to name)
	description string
	// set of fixture names and args to associate with the testcase's
	// before-run stage
	before map[string][]string
	// set of tests that are run as part of this file
	units []Runnable
}

// Context returns the file's Context pointer
func (f *file) Context() *Context {
	return f.ctx
}

// Append appends a runnable test element to the file
func (f *file) Append(r Runnable) {
	f.units = append(f.units, r)
}

// Run executes the tests in the file
func (f *file) Run(t *testing.T) {
	if f.ctx.Fixtures != nil {
		for fname := range f.before {
			f := f.ctx.Fixtures.Get(fname)
			if f == nil {
				t.Fatalf("failed to find required fixture %s", fname)
			}
			f.Start()
			defer f.Stop()
		}
	}
	t.Run(f.path, func(t *testing.T) {
		for _, unit := range f.units {
			unit.Run(t)
		}
	})
}
