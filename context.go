// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

// Context contains fixtures and other state that is passed to the Parse
// function
type Context struct {
	// Fixtures is a pointer to the fixture registry used by the test files
	Fixtures FixtureRegistry
}

func NewContext() *Context {
	return &Context{
		Fixtures: Fixtures,
	}
}
