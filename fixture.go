// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"strings"
)

var (
	// Fixtures is the fixture registry for gdt
	Fixtures = &fixReg
	fixReg   = fixtureRegistry{
		entries: map[string]Fixture{},
	}
)

// A Fixture allows state to be passed from setups
type Fixture interface {
	// Start sets up the fixture
	Start()
	// Stop tears down the fixture, cleaning up any owned resources
	Stop()
	// HasState returns true if the fixture contains some state with the given
	// key
	HasState(string) bool
	// State returns the state data at the given key, or nil if no such state
	// key is managed by the fixture
	State(string) interface{}
}

// FixtureRegistry describes something that can register and return fixtures
type FixtureRegistry interface {
	// Register associates a fixture to a fixture name
	Register(string, Fixture)
	// Get returns a fixture with a given name or nil if no such fixture was
	// recognized in the fixture registry
	Get(string) Fixture
	// List returns a slice of fixtures registered with the registry
	List() []Fixture
}

// fixtureRegistry stores all known Fixtures
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
	if fr == nil {
		return nil
	}
	return fr.entries[strings.ToLower(name)]
}

// List returns a slice of registered fixtures
func (fr *fixtureRegistry) List() []Fixture {
	res := make([]Fixture, len(fr.entries))
	x := 0
	for _, f := range fr.entries {
		res[x] = f
		x++
	}
	return res
}
