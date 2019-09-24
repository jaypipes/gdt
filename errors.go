package gdt

import "errors"

var (
	// ErrNilTestcase is returned when attempting to call a Testcase method
	// with a nil Testcase pointer
	ErrNilTestcase = errors.New("nil test file")
	// ErrUnknownParser is returned when a testcase type name string isn't
	// found in the list of registered test case parsers
	ErrUnknownParser = errors.New("unknown parser")
	// ErrInvalid is returned when a test file is parsed however the test
	// case(s) contained in the test file were somehow not valid
	ErrInvalid = errors.New("invalid test")
)
