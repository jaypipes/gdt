// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"context"

	gdtcontext "github.com/jaypipes/gdt-core/context"
)

// NewContext returns a new `context.Context` that can be passed to a
// `Runnable` (a `Suite` or `Scenario` returned from the `From` function).
//
// If no set of plugins are supplied as a modifier in this function, the
// returned context will have the default set of gdt plugins registered in it.
func NewContext(mods ...gdtcontext.ContextModifier) context.Context {
	ctx := gdtcontext.New(mods...)
	plugins := gdtcontext.Plugins(ctx)
	if len(plugins) == 0 {
		for _, p := range knownPlugins.List() {
			ctx = gdtcontext.RegisterPlugin(ctx, p)
		}
	}
	return ctx
}
