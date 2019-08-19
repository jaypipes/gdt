package gdt

import (
	gdterrors "github.com/jaypipes/gdt/errors"
	"github.com/jaypipes/gdt/interfaces"
	"github.com/jaypipes/gdt/testcase"
)

// FromFile returns a Testcase after reading a supplied filepath and parsing
// the file
func FromFile(fp string) (interfaces.Testcase, error) {
	tc, contents, err := testcase.New().From(fp)
	if err != nil {
		return nil, err
	}
	parser, found := parsers[tc.Type()]
	if !found {
		return nil, gdterrors.ErrUnknownParser
	}
	err = parser.Parse(tc, contents)
	if err != nil {
		return nil, err
	}
	return tc, nil
}
