package suite

import (
	"testing"

	"github.com/jaypipes/gdt/interfaces"
)

// suite contains Testcases, one for each YAML file in a given directory
type suite struct {
	t *testing.T
	// the fixture registry used by the test cases
	fr interfaces.FixtureRegistry
	// Filepath is the filepath to the test suite directory
	filepath string
	// Name for the test suite (defaults to Filepath
	name string
	// Description of the test suite (defaults to Name)
	description string
	// Collection of test cases in this suite
	testcases []interfaces.Runnable
}

func (s *suite) Run() {
	for _, t := range s.testcases {
		t.Run()
	}
}

// New returns a new `Suite`. The function accepts a pointer to the testing.T
// that represents the golang testing framework and accepts zero or more
// `WithOption` values that affect the returned test suite.
//
// Usage:
//
// import "github.com/jaypipes/gdt/suite"
//
// func TestBooksAPI(t *testing.T) {
//     tc := suite.New(t, suite.WithName("books_api"))
//     t.Run()
// }
func New(t *testing.T, opts ...WithOption) interfaces.Runnable {
	useOpts := mergeOptions(opts...)
	tc := &suite{t: t}
	if useOpts.Description != "" {
		tc.description = useOpts.Description
	}
	if useOpts.Name != "" {
		tc.name = useOpts.Name
	}
	if useOpts.FixtureRegistry != nil {
		tc.fr = useOpts.FixtureRegistry
	}
	return tc
}

// FromDir returns a gdt.Suite containing gdt.Testcase objects, one
// for each YAML file in the given directory
func (s *Suite) From(fp string) (interfaces.Runnable, error) {
	return &Suite{
		filepath: fp,
	}
}
