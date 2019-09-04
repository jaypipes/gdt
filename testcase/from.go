package testcase

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"

	gdterrors "github.com/jaypipes/gdt/errors"
	"github.com/jaypipes/gdt/interfaces"
)

type testcaseSpec struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Require     []string `json:"require"`
}

// From parses a gdt YAML file and populates the Testcase with appropriate
// attributes
func (tc *testcase) From(
	fp string,
) (interfaces.Testcase, []byte, error) {
	if tc == nil {
		return nil, nil, gdterrors.ErrNilTestcase
	}
	// We do a double-parse of the test file. The first pass determines the
	// type of test by simply looking for a "type" top-level element in the
	// YAML. If no "type" element was found, the test type defaults to HTTP.
	// Once the type is determined, then the test case module (e.g. gdt/http)
	// is called to parse the file into the case type-specific schema
	f, err := os.Open(fp)
	if err != nil {
		return nil, nil, err
	}
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, nil, err
	}
	tcs := testcaseSpec{}
	if err = yaml.Unmarshal(contents, &tcs); err != nil {
		return nil, nil, err
	}

	// The Testcase may have already had its attributes set using WithOption.
	// Those values are overrides and should not be replaced by any value read
	// from the test file
	if tc.name == "" && tcs.Name != "" {
		tc.name = tcs.Name
	}
	if tc.description == "" && tcs.Description != "" {
		tc.description = tcs.Description
	}

	tc.typ = strings.ToLower(tcs.Type)

	if len(tcs.Require) > 0 {
		tc.before = make(map[string][]string, len(tcs.Require))
		// TODO(jaypipes): Parse a function-call interface from string...
		for _, elem := range tcs.Require {
			tc.before[elem] = []string{}
		}
	}

	return tc, contents, nil
}
