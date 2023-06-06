// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import "context"

type ContextKey string

const (
	fixturesContextKey = ContextKey("gdt.fixtures")
)

// GetFixturesFromContext returns the set of gdt Fixtures from the supplied
// context
func GetFixturesFromContext(ctx context.Context) FixtureRegistry {
	if v := ctx.Value(fixturesContextKey); v != nil {
		return v.(FixtureRegistry)
	}
	return nil
}

// NewContext returns a new gdt testing context that can be passed to a
// Runnable
func NewContext() context.Context {
	return context.WithValue(context.Background(), fixturesContextKey, Fixtures)
}
