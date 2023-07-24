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
	// WithDebug informs gdt to output extra debugging information. You can
	// supply zero or more `io.Writer` objects to the function.
	//
	// If no `io.Writer` objects are supplied, gdt will output debug messages
	// using the `testing.T.Log[f]()` function. This means that you will only
	// get these debug messages if you call the `go test` tool with the `-v`
	// option (either as `go test -v` or with `go test -v=test2json`.
	//
	// ```go
	//
	//	func TestExample(t *testing.T) {
	//		   require := require.New(t)
	//		   fp := filepath.Join("testdata", "example.yaml")
	//		   f, err := os.Open(fp)
	//		   require.Nil(err)
	//
	//		   ctx := gdtcontext.New(gdtcontext.WithDebug())
	//	    s, err := scenario.FromReader(
	//	        f,
	//	        scenario.WithPath(fp),
	//	        scenario.WithContext(ctx),
	//	    )
	//		   require.Nil(err)
	//		   require.NotNil(s)
	//
	//		   err = s.Run(ctx, t)
	//		   require.Nil(err)
	//	}
	//
	// ```
	//
	// If you want gdt to log extra debugging information about tests and
	// assertions to a different file or collecting buffer, pass it a context
	// with a debug `io.Writer`:
	//
	// ```go
	// f := ioutil.TempFile("", "mytest*.log")
	// ctx := gdtcontext.New(gdtcontext.WithDebug(f))
	// ```
	//
	// ```go
	// var b bytes.Buffer
	// w := bufio.NewWriter(&b)
	// ctx := gdtcontext.New(gdtcontext.WithDebug(w))
	// ```
	//
	// you can then inspect the debug "log" and do whatever you'd like with it.
	WithDebug = gdtcontext.WithDebug
	// SetDebug sets gdt's debug logging to the supplied `io.Writer`.
	//
	// The `writers` parameters is optional. If no `io.Writer` objects are
	// supplied, gdt will output debug messages using the `testing.T.Log[f]()`
	// function. This means that you will only get these debug messages if you
	// call the `go test` tool with the `-v` option (either as `go test -v` or
	// with `go test -v=test2json`.
	SetDebug = gdtcontext.SetDebug
	// NewJSONFixture takes a string, some bytes or an io.Reader and returns a
	// new gdttypes.Fixture that can have its state queried via JSONPath
	NewJSONFixture = jsonfix.New
)
