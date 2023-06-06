// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package http

import (
	"context"
	"errors"
	"testing"

	"github.com/jaypipes/gdt"
)

var (
	ErrExpectedLocationHeader = errors.New("Expected Location HTTP Header in previous response")
)

type TestCaseDefaults struct {
	BaseURL string `json:"base_url,omitempty"`
}

// BaseURLFromContext returns the base URL to use when constructing HTTP
// requests. If the Defaults is non-nil and has a BaseURL value, use that.
// Otherwise we look up a base URL from the context's fixtures.
func (d *TestCaseDefaults) BaseURLFromContext(ctx context.Context) string {
	// If the httpFile has been manually configured and the configuration
	// contains a base URL, use that. Otherwise, check to see if there is a
	// fixture in the registry that has an "http.base_url" state key and use
	// that if found.
	if d != nil && d.BaseURL != "" {
		return d.BaseURL
	}
	// query the fixture registry to determine if any of them contain an
	// http.base_url state attribute.
	for _, f := range gdt.GetFixturesFromContext(ctx).List() {
		if f.HasState(StateKeyBaseURL) {
			return f.State(StateKeyBaseURL).(string)
		}
	}
	return ""
}

// TestCase describes a set of tests of HTTP requests and responses
type TestCase struct {
	// Defaults contains the default configuration options for this test case
	Defaults *TestCaseDefaults `json:"defaults,omitempty"`
	// Specs contains one or more specifications for a single HTTP
	// request/response test
	Specs []*TestSpec `json:"tests,omitempty"`
}

// Run executes all tests described by the test case.
func (c *TestCase) Run(ctx context.Context, t *testing.T) {
	for _, s := range c.Specs {
		ctx = s.Run(ctx, t)
	}
}
