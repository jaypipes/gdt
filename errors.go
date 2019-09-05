package gdt

import "errors"

var (
	// ErrNilTestcase is returned when attempting to call a Testcase method
	// with a nil Testcase pointer
	ErrNilTestcase = errors.New("nil testcase")
	// ErrUnknownParser is returned when a testcase type name string isn't
	// found in the list of registered test case parsers
	ErrUnknownParser = errors.New("unknown parser")
)
