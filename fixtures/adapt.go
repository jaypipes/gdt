package fixtures

import "github.com/jaypipes/gdt"

type simpleFixture struct {
	starter func()
	stopper func()
	state   map[string]string
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
		if _, ok := f.state[key]; ok {
			return true
		}
	}
	return false
}

// GetState returns a string attribute from the fixture's state map if the
// supplied key exists, otherwise returns empty string
func (f *simpleFixture) State(key string) interface{} {
	if f.state != nil {
		return f.state[key]
	}
	return ""
}

type WithOption struct {
	Starter func()
	Stopper func()
	State   map[string]string
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
func WithState(state map[string]string) WithOption {
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
