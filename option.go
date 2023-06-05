// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

// Option contains configuration information
type Option struct {
	// name is the short name for the test suite. If empty, defaults to Path
	name *string
	// description is a description of the tests contained in the test suite
	description *string
	// path is the filepath to the test suite directory
	path *string
	// typ is the test case type
	typ *string
}

// WithName returns an Option with a name property. When manually creating
// TestCase or TestSuite objects, you can pass WithName("myname") to
// NewTestCase() or NewTestSuite()
func WithName(name string) *Option {
	return &Option{name: &name}
}

// WithDescription returns an Option with a description property. When manually
// creating TestCase or TestSuite objects, you can pass
// WithDescription("mydescription") to NewTestCase() or NewTestSuite()
func WithDescription(description string) *Option {
	return &Option{description: &description}
}

// WithPath returns an Option with a path property. When manually creating
// TestCase or TestSuite objects, you can pass WithPath("mypath") to
// NewTestCase() or NewTestSuite()
func WithPath(path string) *Option {
	return &Option{path: &path}
}

// WithType returns an Option with a name property. When manually creating
// TestCase or TestSuite objects, you can pass WithType("myname") to
// NewTestCase() or NewTestSuite()
func WithType(typ string) *Option {
	return &Option{typ: &typ}
}

// mergeOptions merges all supplied options into a single one
func mergeOptions(opts []*Option) *Option {
	merged := &Option{}
	for _, opt := range opts {
		if opt.name != nil {
			merged.name = opt.name
		}
		if opt.description != nil {
			merged.description = opt.description
		}
		if opt.path != nil {
			merged.path = opt.path
		}
		if opt.path != nil {
			merged.typ = opt.typ
		}
	}
	return merged
}
