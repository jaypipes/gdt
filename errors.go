// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import "errors"

var (
	// ErrUnknownParser is returned when a testcase type name string isn't
	// found in the list of registered test case parsers
	ErrUnknownParser = errors.New("unknown parser")
	// ErrInvalidYAML is returned when a test file could not be parsed because
	// the file contents were not valid YAML
	ErrInvalidYAML = errors.New("file contents not valid YAML")
	// ErrInvalid is returned when a test file is parsed however the test
	// case(s) contained in the test file were somehow not valid
	ErrInvalid = errors.New("invalid test")
)
