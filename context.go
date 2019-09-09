package gdt

// Context contains fixtures and other state that is passed to the Parse
// function
type Context struct {
	// FixtureRegistry is a pointer to the fixture registry used by the test
	// files
	Fixtures FixtureRegistry
}
