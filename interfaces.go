package gdt

import "testing"

// Runnable represents things that have a simple Run() method that accepts a
// pointer to a testing.T. Example things that implement this interface are
// gdt.TestCase and gdt.TestSuite
type Runnable interface {
	Run(*testing.T)
}

// Appendable simply allows some runnable thing to be added to it
type Appendable interface {
	Append(Runnable)
}

// ContextAppendable is an Appendable that can return a Context
type ContextAppendable interface {
	Appendable
	Context() *Context
}
