package testcase

import (
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"

	gdt_http "./http"
)

type fixtureSpec struct {
	Name string
	Args []string
}

type setupSpec struct {
	Fixtures map[string]*fixtureSpec
}

type testcaseSpec struct {
	Type        string    `json:"type"`
	Name        string    `json:"name",omitempty`
	Description string    `json:"description",omitempty`
	Setup       setupSpec `json:"setup"`
}

// From parses a gdt YAML file and populates the Testcase with appropriate
// attributes
func (tc *Testcase) From(fp string) (*Testcase, error) {
	if tc == nil {
		return nil, ErrNilTestcase
	}
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
	tcs := testcaseSpec{}
	if err = yaml.Unmarshal(contents, &tcs); err != nil {
		return nil, err
	}

	// The Testcase may have already had its attributes set using WithOption.
	// Those values are overrides and should not be replaced by any value read
	// from the test file
	if tc.Name == "" && tcs.Name != "" {
		tc.Name = tcs.Name
	}
	if tc.Description == "" && tcs.Description != "" {
		tc.Description = tcs.Description
	}

	switch tp.Type {
	case "http", "":
		{
			err := gdt_http.Parse(tc, contents)
			if err != nil {
				return nil, err
			}
			return tc, nil
		}
	}
}
