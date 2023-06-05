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
	gdt.Parsers.Register(&fooParser{}, "foo")
	tests := []struct {
		name     string
		contents []byte
		err      error
		tc       *gdt.TestCase
	}{
		{
			name:     "empty content",
			contents: []byte{},
			err:      gdt.ErrUnknownParser,
			tc:       nil,
		},
		{
			name: "bad YAML",
			contents: []byte(`
			bad YAML, Indy
			`),
			err: gdt.ErrInvalidYAML,
			tc:  nil,
		},
		{
			name: "unknown parser",
			contents: []byte(`type: bar
name: bar test
`),
			err: gdt.ErrUnknownParser,
			tc:  nil,
		},
		{
			name: "found parser",
			contents: []byte(`type: foo
name: foo test
`),
			err: nil,
			tc: &gdt.TestCase{
				Type: "foo",
				Name: "foo test",
			},
		},
	}

	for _, test := range tests {
		tc, err := gdt.FromBytes(test.contents)
		if err != nil {
			assert.Equal(t, test.err, err)
		} else {
			assert.Equal(t, test.tc, tc)
		}
	}
}
