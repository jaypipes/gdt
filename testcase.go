// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"strings"
	"testing"
)

var (
	// Parsers is the parser registry for gdt
	Parsers = &parReg
	parReg  = parserRegistry{
		entries: map[string]Parser{},
	}
)

// Parser is the driver interface for parsers of different types of tests
type Parser interface {
	// Parse the supplied raw contents and append any elements to the supplied
	// ContextAppendable
	Parse(ca Appendable, contents []byte) error
}

// parserRegistry stores all known Parsers
type parserRegistry struct {
	entries map[string]Parser
}

// Register associates a parser to one or more types of test (files)
func (pr *parserRegistry) Register(p Parser, testTypes ...string) {
	for _, tt := range testTypes {
		pr.entries[strings.ToLower(tt)] = p
	}
}

// Get returns the parser for a given test type or nil if no such parser was
// recognized in the parser registry
func (pr *parserRegistry) Get(testType string) Parser {
	if pr == nil {
		return nil
	}
	return pr.entries[strings.ToLower(testType)]
}

// List returns a slice of registered parsers
func (pr *parserRegistry) List() []Parser {
	res := make([]Parser, len(pr.entries))
	x := 0
	for _, p := range pr.entries {
		res[x] = p
		x++
	}
	return res
}

// TestCase is a generalized gdt test case file. It contains a set of Runnable
// test units.
type TestCase struct {
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
	// path is the filepath to the test case
	path string `json:"-"`
}

// Append appends a runnable test element to the test case
func (tc *TestCase) Append(r Runnable) {
	tc.units = append(tc.units, r)
}

// Run executes the tests in the test case
func (tc *TestCase) Run(t *testing.T, ctx *Context) {
	if ctx.Fixtures != nil {
		for _, fname := range tc.Require {
			fix := ctx.Fixtures.Get(fname)
			if fix == nil {
				t.Fatalf("failed to find required fixture '%s'", fname)
			}
			V2("file.file:Run", "starting fixture %s\n", fname)
			fix.Start()
			defer fix.Stop()
		}
	}
	t.Run(tc.path, func(t *testing.T) {
		for _, unit := range tc.units {
			unit.Run(t, ctx)
		}
	})
}
