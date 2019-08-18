package testcase

// Assertion is used by the test framework to generalize a thing that asserts
// some truth value
type Assertion interface {
	// Assert returns whether the assertion succeeded and a string indicating
	// why the assertion did not succeed
	Assert() (bool, string)
}
