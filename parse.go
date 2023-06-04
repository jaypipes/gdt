// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
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
	Parse(ca ContextAppendable, contents []byte) error
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

type fileSchema struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Require     []string `json:"require"`
}

// parse parses the supplied file and returns an error if any syntax or
// validation failed
//
// We do a double-parse of the test file. The first pass determines the
// type of test by simply looking for a "type" top-level element in the
// YAML. If no "type" element was found, the test type defaults to HTTP.
// Once the type is determined, then the test case module (e.g. gdt/http)
// is called to parse the file into the case type-specific schema
func (tf *file) parse(r io.Reader) error {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return parseBytes(tf, contents)
}

func parseBytes(
	tf *file,
	contents []byte,
) error {
	tfs := fileSchema{}
	if err := yaml.Unmarshal(contents, &tfs); err != nil {
		return ErrInvalidYAML
	}
	parser := Parsers.Get(tfs.Type)
	if parser == nil {
		return ErrUnknownParser
	}

	tf.typ = strings.ToLower(tfs.Type)
	tf.name = tfs.Name
	tf.description = tfs.Description
	tf.require = tfs.Require

	if err := parser.Parse(tf, contents); err != nil {
		return err
	}

	return nil
}
