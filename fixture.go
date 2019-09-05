package gdt

import (
	"strings"
)

var (
	// Fixtures is the default fixture registry for gdt
	Fixtures      = &gdtFixtureReg
	gdtFixtureReg = fixtureRegistry{
		entries: map[string]Fixture{},
	}
)

// implements interfaces.FixtureRegistry
type fixtureRegistry struct {
	entries map[string]Fixture
}

// Register associates a fixture to a fixture name
func (fr *fixtureRegistry) Register(name string, f Fixture) {
	fr.entries[strings.ToLower(name)] = f
}

// Get returns a fixture with a given name or nil if no such fixture was
// recognized in the fixture registry
func (fr *fixtureRegistry) Get(name string) Fixture {
	return fr.entries[strings.ToLower(name)]
}

// List returns an array of fixtures
func (fr *fixtureRegistry) List() []Fixture {
	res := make([]Fixture, len(fr.entries))
	x := 0
	for _, f := range fr.entries {
		res[x] = f
		x++
	}
	return res
}
