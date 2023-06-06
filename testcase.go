// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"context"
	"io"
	"io/ioutil"
	"path"
	gopath "path"
	"strings"
	"testing"

	"github.com/ghodss/yaml"
)

// TestCase is a generalized gdt test case file. It contains a set of Runnable
// test units.
type TestCase struct {
	// Path is the filepath to the test case
	Path string `json:"-"`
	// Type is the type of test contained in the file. Defaults to "http"
	Type string `json:"type,omitempty"`
	// Name is the short name for the test case. If empty, defaults to Path
	Name string `json:"name,omitempty"`
	// Description is a description of the tests contained in the test case
	Description string `json:"description,omitempty"`
	// Require specifies an ordered list of fixtures the test case depends on
	Require []string `json:"require"`
	// set of tests that are run as part of this file
	units []Runnable `json:"-"`
}

// Append appends a runnable test element to the test case
func (tc *TestCase) Append(r Runnable) {
	tc.units = append(tc.units, r)
}

// Run executes the tests in the test case
func (tc *TestCase) Run(ctx context.Context, t *testing.T) {
	fixtures := GetFixturesFromContext(ctx)
	if fixtures != nil {
		for _, fname := range tc.Require {
			fix := fixtures.Get(fname)
			if fix == nil {
				t.Fatalf("failed to find required fixture '%s'", fname)
			}
			V2("file.file:Run", "starting fixture %s\n", fname)
			fix.Start()
			defer fix.Stop()
		}
	}
	t.Run(tc.Name, func(t *testing.T) {
		for _, unit := range tc.units {
			unit.Run(ctx, t)
		}
	})
}

// NewTestCase returns a new TestCase
func NewTestCase(options ...*Option) *TestCase {
	merged := mergeOptions(options)
	tc := &TestCase{}
	if merged.path != nil {
		tc.Path = *merged.path
	}
	if merged.name != nil {
		tc.Name = *merged.name
	}
	if merged.description != nil {
		tc.Description = *merged.description
	}
	if merged.typ != nil {
		tc.Type = *merged.typ
	}

	// default the name of the test case to the basename of the file path
	if tc.Name == "" && tc.Path != "" {
		tc.Name = path.Base(tc.Path)
	}
	return tc
}

// NewTestCaseFromReader parses the supplied io.Reader and returns a TestCase
// representing the contents in the reader. Returns an error if any syntax or
// validation failed
func NewTestCaseFromReader(r io.Reader, path string) (Runnable, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewTestCaseFromBytes(contents, path)
}

// NewTestCaseFromBytes returns a TestCase after parsing the supplied contents
func NewTestCaseFromBytes(contents []byte, path string) (*TestCase, error) {
	// We do a double-parse of the test file. The first pass determines the
	// type of test by simply looking for a "type" top-level element in the
	// YAML. If no "type" element was found, the test type defaults to HTTP.
	// Once the type is determined, then the test case module (e.g. gdt/http)
	// is called to parse the file into the case type-specific schema
	tc := TestCase{}
	if err := yaml.Unmarshal(contents, &tc); err != nil {
		return nil, ErrInvalidYAML
	}

	tc.Path = path

	// default the name of the test case to the basename of the file path
	if tc.Name == "" {
		tc.Name = gopath.Base(path)
	}

	tc.Type = strings.ToLower(tc.Type)
	parser := TestTypeParsers.Get(tc.Type)
	if parser == nil {
		return nil, ErrUnknownParser
	}

	if err := parser.Parse(&tc, contents); err != nil {
		return nil, err
	}

	return &tc, nil
}
