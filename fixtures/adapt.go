// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package fixtures

import (
	"strings"

	"github.com/jaypipes/gdt"
)

type simpleFixture struct {
	starter func()
	stopper func()
	state   map[string]interface{}
}

// Start sets up any resources the fixture uses
func (f *simpleFixture) Start() {
	if f.starter != nil {
		f.starter()
	}
}

// Stop cleans up any resources the fixture uses
func (f *simpleFixture) Stop() {
	if f.stopper != nil {
		f.stopper()
	}
}

// HasState returns true if the fixture has a state attribute with the supplied
// key
func (f *simpleFixture) HasState(key string) bool {
	if f.state != nil {
		_, ok := f.state[strings.ToLower(key)]
		return ok
	}
	return false
}

// State returns a piece of state from the fixture's state map if the supplied
// key exists, otherwise returns nil
func (f *simpleFixture) State(key string) interface{} {
	if f.state != nil {
		return f.state[strings.ToLower(key)]
	}
	return nil
}

type WithOption struct {
	Starter func()
	Stopper func()
	State   map[string]interface{}
}

// WithStart allows a starter functor to be adapted into a fixture
func WithStart(starter func()) WithOption {
	return WithOption{Starter: starter}
}

// WithStop allows a stopper functor to be adapted into a fixture
func WithStop(stopper func()) WithOption {
	return WithOption{Stopper: stopper}
}

// WithState allows a map of state key/values to be adapted into a fixture
func WithState(state map[string]interface{}) WithOption {
	return WithOption{State: state}
}

// Adapt returns a simple object that implements the interfaces.Fixture
// interface from one or more WithOptions describing starter, stopper functions
// or a state map
func Adapt(opts ...*WithOption) gdt.Fixture {
	if len(opts) == 0 {
		panic("gdt.fixtures.Adapt should be called with at least one WithOption")
	}
	res := &simpleFixture{}

	for _, opt := range opts {
		if opt.Starter != nil {
			res.starter = opt.Starter
		}
		if opt.Stopper != nil {
			res.stopper = opt.Stopper
		}
		if opt.State != nil {
			res.state = opt.State
		}
	}
	return res
}
