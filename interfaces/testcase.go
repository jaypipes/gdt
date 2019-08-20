package interfaces

import "testing"

// Testcase describes some related test units
type Testcase interface {
	Typed
	Filepath
	Named
	T() *testing.T
}
