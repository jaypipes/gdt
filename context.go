package gdt

// Context contains fixtures and other state that is passed to the Parse
// function
type Context struct {
	// Fixtures is a pointer to the fixture registry used by the test files
	Fixtures FixtureRegistry
	// Basedir is an absolute filepath to the directory that contains the test
	// files
	Basedir string
}
