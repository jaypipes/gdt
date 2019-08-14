package testcase

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ghodss/yaml"

	gdt_http "github.com/jaypipes/gdt/http"
)

type testCaseTypeProber struct {
	Type string `json:"type"`
}

// FromFile reads a GDT test from the supplied filepath and returns a
// gdt.TestCase describing the test
func FromFile(fp string, opts ...gdt.WithOption) *gdt.TestCase, error {
	// We do a double-parse of the test file. The first pass determines the
	// type of test by simply looking for a "type" top-level element in the
	// YAML. If no "type" element was found, the test type defaults to HTTP.
	// Once the type is determined, then the test case module (e.g. gdt/http)
	// is called to parse the file into the case type-specific schema
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	tp := testTypeProber{}
	if err = yaml.Unmarshal(contents, &tp); err != nil {
		return nil, err
	}

	switch tp.Type {
	case "http", "":
		{
			tc, err := gdt_http.NewTestCaseFromYAML(contents)
			if err != nil {
				return nil, err
			}
			return tc, nil
		}
	}
}
