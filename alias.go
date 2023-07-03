// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	gdtcontext "github.com/jaypipes/gdt-core/context"
	jsonfix "github.com/jaypipes/gdt-core/fixture/json"
)

var (
	// RegisterFixture registers a named fixtures with the context
	RegisterFixture = gdtcontext.RegisterFixture
	// WithFixtures sets a context's Fixtures
	WithFixtures = gdtcontext.WithFixtures
	// WithDebug sets a context's Debug writer. If you want gdt to log extra
	// debugging information about tests and assertions, pass it a context with
	// a debug writer:
	//
	// ```go
	// f := ioutil.TempFile("", "mytest*.log")
	// ctx := gdt.NewContext(gdt.WithDebug(f))
	// ```
	//
	// you can then inspect the debug "log" and do whatever you'd like with it.
	//
	// Or you could pass a console writer and just have gdt write to the
	// console its debugging information:
	//
	// ```go
	// ctx := gdt.NewContext(gdt.WithDebug(os.Stdout))
	// ```
	WithDebug = gdtcontext.WithDebug
	// SetDebug sets gdt's debug logging to the supplied `io.Writer`
	SetDebug = gdtcontext.SetDebug
	// NewJSONFixture takes a string, some bytes or an io.Reader and returns a
	// new gdttypes.Fixture that can have its state queried via JSONPath
	NewJSONFixture = jsonfix.New
)
