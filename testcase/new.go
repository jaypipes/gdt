package testcase

// NewTestCase returns a new `gdt.TestCase` for an HTTP test case. The function
// accepts zero or more `gdt.WithOption` values that affect the returned test
// case.
//
// Usage:
//
//   tc := gdt.NewTestCase(gdt.WithName("books_api"))
func New(opts ...WithOption) *gdt.TestCase {
	useOpts := mergeOptions(opts...)
	t := &TestCase{
		TestCaseType: TestCaseTypeHTTP,
	}
	if useOpts.Name != "" {
		t.Name = useOpts.Name
	} else {
		// Default the test case name to the filepath
	}
	if useOpts.Description != "" {
		t.Description = useOpts.Description
	}
	if useOpts.Filepath != "" {
		t.Filepath = useOpts.Filepath
	}
	return t
}
