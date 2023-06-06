// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt_test

import (
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/stretchr/testify/assert"
)

type fooParser struct{}

func (p *fooParser) Parse(gdt.Appendable, []byte) error {
	return nil
}

func TestFromBytes(t *testing.T) {
	gdt.TestTypeParsers.Register(&fooParser{}, "foo")
	tests := []struct {
		name     string
		contents []byte
		err      error
		path     string
		tc       *gdt.TestCase
	}{
		{
			name:     "empty content",
			contents: []byte{},
			err:      gdt.ErrUnknownParser,
			path:     "/my/path/to/file.yaml",
			tc:       nil,
		},
		{
			name: "bad YAML",
			contents: []byte(`
			bad YAML, Indy
			`),
			err:  gdt.ErrInvalidYAML,
			path: "/my/path/to/file.yaml",
			tc:   nil,
		},
		{
			name: "unknown parser",
			contents: []byte(`type: bar
name: bar test
`),
			err:  gdt.ErrUnknownParser,
			path: "/my/path/to/file.yaml",
			tc:   nil,
		},
		{
			name: "good parse, name included",
			contents: []byte(`type: foo
name: foo test
`),
			err:  nil,
			path: "/my/path/to/file.yaml",
			tc: &gdt.TestCase{
				Path: "/my/path/to/file.yaml",
				Type: "foo",
				Name: "foo test",
			},
		},
		{
			name: "good parse, name as base path",
			contents: []byte(`type: foo
`),
			err:  nil,
			path: "/my/path/to/file.yaml",
			tc: &gdt.TestCase{
				Path: "/my/path/to/file.yaml",
				Type: "foo",
				Name: "file.yaml",
			},
		},
	}

	for _, test := range tests {
		tc, err := gdt.NewTestCaseFromBytes(test.contents, test.path)
		if err != nil {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.tc, tc)
		}
	}
}
