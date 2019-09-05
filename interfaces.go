package gdt

import "testing"

// A Fixture allows state to be passed from setups
type Fixture interface {
	Start()
	Stop()
	HasState(string) bool
	State(string) string
}

// FixtureRegistry describes something that can register and return fixtures
type FixtureRegistry interface {
	Register(string, Fixture)
	Get(string) Fixture
	List() []Fixture
}

// Parser is the driver interface for parsers of different types of tests
type Parser interface {
	// Parse the supplied raw contents and append any elements to the supplied
	// Appendable
	Parse(tf *TestFile, contents []byte) error
}

// Runnable represents things that have a simple Run() method that accepts a
// pointer to a testing.T. Example things that implement this interface are
// gdt.TestCase and gdt.TestSuite
type Runnable interface {
	Run(*testing.T)
}

// Appendable simply allows some runnable thing to be added to it
type Appendable interface {
	Append(Runnable)
}
