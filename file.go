package gdt

import (
	"testing"
)

// TestFile describes the tests in a single gdt test file. Implements
// interfaces.Appendable and interfaces.Runnable
type TestFile struct {
	ctx *context
	// typ is the type of test case
	typ string
	// path is the filepath to the test file
	path string
	// name for the overall test
	name string
	// description of the test (defaults to name)
	description string
	// set of fixture names and args to associate with the testcase's
	// before-run stage
	before map[string][]string
	// set of tests that are run as part of this test case
	units []Runnable
}

// Fixtures returns a pointer to the fixture registry
func (tf *TestFile) Fixtures() []Fixture {
	return tf.ctx.fr.List()
}

// Append appends a runnable test element to the test file
func (tf *TestFile) Append(r Runnable) {
	tf.units = append(tf.units, r)
}

// Run executes the tests in the test case
func (tf *TestFile) Run(t *testing.T) {
	if tf.ctx.fr != nil {
		for fname := range tf.before {
			f := tf.ctx.fr.Get(fname)
			if f == nil {
				t.Fatalf("failed to find required fixture %s", fname)
			}
			f.Start()
			defer f.Stop()
		}
	}
	t.Run(tf.name, func(t *testing.T) {
		for _, test := range tf.units {
			test.Run(t)
		}
	})
}
