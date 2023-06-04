// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fooParser struct{}

func (p *fooParser) Parse(ContextAppendable, []byte) error {
	return nil
}

func TestParseBytes(t *testing.T) {
	Parsers.Register(&fooParser{}, "foo")
	tests := []struct {
		name     string
		contents []byte
		exp      error
	}{
		{
			name:     "empty content",
			contents: []byte{},
			exp:      ErrUnknownParser,
		},
		{
			name: "bad YAML",
			contents: []byte(`
			bad YAML, Indy
			`),
			exp: ErrInvalidYAML,
		},
		{
			name: "unknown parser",
			contents: []byte(`type: bar
name: bar test
`),
			exp: ErrUnknownParser,
		},
		{
			name: "found parser",
			contents: []byte(`type: foo
name: foo test
`),
			exp: nil,
		},
	}

	for _, test := range tests {
		tf := file{
			name: test.name,
		}
		got := parseBytes(&tf, test.contents)
		assert.Equal(t, test.exp, got)
	}
}
