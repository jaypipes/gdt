package interfaces

// Testcase describes some related test units
type Testcase interface {
	Typed
	Filepath
	Runnable
	// AppendRunnable adds a Runnable thing to be run when the Testable executed
	AppendRunnable(Runnable)
}
