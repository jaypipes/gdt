// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"context"
	"os"
	"path/filepath"
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
func (s *TestSuite) Run(ctx context.Context, t *testing.T) {
	for _, unit := range s.units {
		unit.Run(ctx, t)
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

// NewTestSuiteFromDir reads the supplied directory path and returns a
// TestSuite representing the suite of test cases in that directory.
func NewTestSuiteFromDir(dirPath string) (*TestSuite, error) {
	// List YAML files in the directory and parse each into a testable unit
	s := NewTestSuite(WithPath(dirPath))

	if err := filepath.Walk(
		dirPath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			suffix := filepath.Ext(path)
			if suffix != ".yaml" {
				return nil
			}
			f, err := os.Open(path)

			if err != nil {
				return err
			}
			defer f.Close()

			tc, err := NewTestCaseFromReader(f, path)
			if err != nil {
				return err
			}
			s.Append(tc)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return s, nil
}
