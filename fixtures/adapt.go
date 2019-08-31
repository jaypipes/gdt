package fixtures

import "github.com/jaypipes/gdt/interfaces"

type simpleFixture struct {
	starter func()
	cleaner func()
}

// Start sets up any resources the fixture uses
func (f *simpleFixture) Start() {
	if f.starter != nil {
		f.starter()
	}
}

// Cleanup cleans up any resources the fixture uses
func (f *simpleFixture) Cleanup() {
	if f.start != nil {
		f.cleaner()
	}
}

// AdaptStart returns a simple object that implements the interfaces.Fixture
// interface from a function that accepts no arguments and returns no arguments
// and will be run when the fixture is started
func AdaptStart(start func()) interfaces.Fixture {
	return &simpleFixture{starter: start}
}
