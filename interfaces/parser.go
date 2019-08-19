package interfaces

// Parser is the driver interface for parsers of different types of tests
type Parser interface {
	// Parse modifies the Testcase after parsing the supplied test case raw
	// contents
	Parse(tc Testcase, contents []byte) error
}
