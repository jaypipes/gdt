package interfaces

import "io"

// RunResult is what is returned from Runnable.Run
type RunResult interface {
	OK() bool
	Skipped() bool
	Errors() []error
}

// Runnable describes the interface that a Testcase runs
type Runnable interface {
	Named
	// Run takes two io.Writers, one for the normal output stream, the other
	// for the error stream, and returns a RunResult
	Run(ow, ew io.Writer) RunResult
}
