package gdt

import (
	"strings"

	"github.com/jaypipes/gdt/interfaces"
)

var (
	parsers  = map[string]interfaces.Parser{}
	fixtures = map[string]interfaces.Fixture{}
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

// RegisterFixture registers a fixture for a fixture names
func RegisterFixture(
	fixture interfaces.Fixture,
	name string,
) {
	fixtures[strings.ToLower(name)] = fixture
}
