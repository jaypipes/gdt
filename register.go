package gdt

import (
	"strings"

	"github.com/jaypipes/gdt/interfaces"
)

var (
	parsers = map[string]interfaces.Parser{}
)

// RegisterParser registers a parser for one or more test case type names
func RegisterParser(
	parser interfaces.Parser,
	types ...string,
) {
	for _, typ := range types {
		parsers[strings.ToLower(typ)] = parser
	}
}
