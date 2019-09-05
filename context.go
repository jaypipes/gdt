package gdt

// context contains fixtures and other state that is passed to the Parse
// function
type context struct {
	// the fixture registry used by the test cases
	fr FixtureRegistry
}
