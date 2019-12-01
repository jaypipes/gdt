package fixtures

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/PaesslerAG/jsonpath"
	"github.com/jaypipes/gdt"
)

type jsonFixture struct {
	data interface{}
}

func (f *jsonFixture) Start() {}

func (f *jsonFixture) Stop() {}

// HasState returns true if the supplied JSONPath expression results in a found
// value in the fixture's data
func (f *jsonFixture) HasState(path string) bool {
	if f.data == nil {
		return false
	}
	got, err := jsonpath.Get(path, f.data)
	if err != nil {
		return false
	}
	if got == nil {
		return false
	}
	return true
}

// GetState returns the value at supplied JSONPath expression
func (f *jsonFixture) State(path string) interface{} {
	if f.data == nil {
		return ""
	}
	got, err := jsonpath.Get(path, f.data)
	if err != nil {
		return ""
	}
	switch got.(type) {
	case string:
		return got.(string)
	case float64:
		return strconv.FormatFloat(got.(float64), 'f', 0, 64)
	default:
		return ""
	}
}

// NewJSONFixture takes an io.Reader and returns a new gdt.Fixture that can
// have its data queried via JSONPath
func NewJSONFixture(r io.Reader) (gdt.Fixture, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	f := jsonFixture{
		data: interface{}(nil),
	}
	if err = json.Unmarshal(b, &f.data); err != nil {
		return nil, err
	}
	return &f, nil
}
