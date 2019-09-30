package gdt

import "errors"

var (
	// ErrUnknownParser is returned when a testcase type name string isn't
	// found in the list of registered test case parsers
	ErrUnknownParser = errors.New("unknown parser")
	// ErrInvalid is returned when a test file is parsed however the test
	// case(s) contained in the test file were somehow not valid
	ErrInvalid = errors.New("invalid test")
)
