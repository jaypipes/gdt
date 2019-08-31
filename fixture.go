package gdt

import (
	"strings"

	"github.com/jaypipes/gdt/interfaces"
)

var (
	// Fixtures is the default fixture registry for gdt
	Fixtures      = &gdtFixtureReg
	gdtFixtureReg = fixtureRegistry{
		entries: map[string]interfaces.Fixture{},
	}
)

// implements interfaces.FixtureRegistry
type fixtureRegistry struct {
	entries map[string]interfaces.Fixture
}

// Register associates a fixture to a fixture name
func (fr *fixtureRegistry) Register(name string, f interfaces.Fixture) {
	fr.entries[strings.ToLower(name)] = f
}

// Get returns a fixture with a given name or nil if no such fixture was
// recognized in the fixture registry
func (fr *fixtureRegistry) Get(name string) interfaces.Fixture {
	return fr.entries[strings.ToLower(name)]
}
