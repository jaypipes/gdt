// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"strings"
)

var (
	// TestTypeParsers is the parser registry for gdt test types
	TestTypeParsers = &ttParReg
	ttParReg        = ttParserRegistry{
		entries: map[string]TestTypeParser{},
	}
)

// TestTypeParser is the driver interface for parsers of different types of tests
type TestTypeParser interface {
	// Parse takes the supplied raw contents and appends any elements to the
	// supplied Appendable
	Parse(Appendable, []byte) error
}

// ttParserRegistry stores all known TestTypeParsers
type ttParserRegistry struct {
	entries map[string]TestTypeParser
}

// Register associates a parser to one or more types of test (files)
func (r *ttParserRegistry) Register(p TestTypeParser, testTypes ...string) {
	for _, tt := range testTypes {
		r.entries[strings.ToLower(tt)] = p
	}
}

// Get returns the parser for a given test type or nil if no such parser was
// recognized in the parser registry
func (r *ttParserRegistry) Get(testType string) TestTypeParser {
	if r == nil {
		return nil
	}
	return r.entries[strings.ToLower(testType)]
}

// List returns a slice of registered parsers
func (r *ttParserRegistry) List() []TestTypeParser {
	res := make([]TestTypeParser, len(r.entries))
	x := 0
	for _, p := range r.entries {
		res[x] = p
		x++
	}
	return res
}
