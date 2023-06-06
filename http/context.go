// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package http

import (
	"context"
	nethttp "net/http"

	"github.com/jaypipes/gdt"
)

const (
	prevResponseContextKey = gdt.ContextKey("gdt.http.prevresp")
)

// getPreviousResponse returns the last HTTP response one of the test units
// executed
func getPreviousResponse(ctx context.Context) *nethttp.Response {
	if v := ctx.Value(prevResponseContextKey); v != nil {
		return v.(*nethttp.Response)
	}
	return nil
}

// storePreviousResponse caches the supplied HTTP response in the supplied
// context and returns the new context
func storePreviousResponse(ctx context.Context, resp *nethttp.Response) context.Context {
	return context.WithValue(ctx, prevResponseContextKey, resp)
}
