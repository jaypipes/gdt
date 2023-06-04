// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package http

import "github.com/jaypipes/gdt"

func init() {
	gdt.Parsers.Register(&httpParser{}, "http", "")
}
