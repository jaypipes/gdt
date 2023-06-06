// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package http

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/jaypipes/gdt"
)

const (
	msgUnsupportedJSONSchemaReference = "unsupported JSONSchema reference URL: %s"
	msgJSONSchemaFileNotFound         = "unable to find JSONSchema file: %s"
)

func init() {
	gdt.TestTypeParsers.Register(&httpParser{}, "http", "")
}

func errUnsupportedJSONSchemaReference(url string) error {
	return fmt.Errorf(msgUnsupportedJSONSchemaReference, url)
}

func errJSONSchemaFileNotFound(path string) error {
	return fmt.Errorf(msgJSONSchemaFileNotFound, path)
}

type httpParser struct{}

// Parse accepts a Testcase and a string of YAML contents from a gdt test file.
// It then parses the HTTP test case and adds the HTTP-specific tests to the
// supplied Testcase
func (p *httpParser) Parse(a gdt.Appendable, contents []byte) error {
	var err error
	tc := TestCase{}
	if err = yaml.Unmarshal(contents, &tc); err != nil {
		return err
	}
	for _, s := range tc.Specs {
		if err = validateResponseAssertion(s.Response); err != nil {
			return err
		}
		if err = validateMethodAndURL(s); err != nil {
			return err
		}
		s.defaults = tc.Defaults
	}
	a.Append(&tc)
	return nil
}

func validateMethodAndURL(s *TestSpec) error {
	if s.URL == "" {
		if s.GET != "" {
			s.Method = "GET"
			s.URL = s.GET
			return nil
		} else if s.POST != "" {
			s.Method = "POST"
			s.URL = s.POST
			return nil
		} else if s.PUT != "" {
			s.Method = "PUT"
			s.URL = s.PUT
			return nil
		} else if s.DELETE != "" {
			s.Method = "DELETE"
			s.URL = s.DELETE
			return nil
		} else if s.PATCH != "" {
			s.Method = "PATCH"
			s.URL = s.PATCH
			return nil
		} else {
			return ErrInvalidAliasOrURL
		}
	}
	if s.Method == "" {
		return ErrInvalidAliasOrURL
	}
	return nil
}

func validateResponseAssertion(resp *ResponseAssertion) error {
	if resp == nil {
		return nil
	}
	if resp.JSON == nil {
		return nil
	}
	if resp.JSON.Schema == "" {
		return nil
	}
	// Ensure any JSONSchema URL specified in response.json.schema exists
	schemaURL := resp.JSON.Schema
	if strings.HasPrefix(schemaURL, "http://") || strings.HasPrefix(schemaURL, "https://") {
		// TODO(jaypipes): Support network lookups?
		return errUnsupportedJSONSchemaReference(schemaURL)
	}
	// Convert relative filepaths to absolute filepaths rooted in the context's
	// testdir after stripping any "file://" scheme prefix
	schemaURL = strings.TrimPrefix(schemaURL, "file://")
	schemaURL, _ = filepath.Abs(schemaURL)

	f, err := os.Open(schemaURL)
	if err != nil {
		return errJSONSchemaFileNotFound(schemaURL)
	}
	defer f.Close()
	if runtime.GOOS == "windows" {
		// Need to do this because of an "optimization" done in the
		// gojsonreference library:
		// https://github.com/xeipuuv/gojsonreference/blob/bd5ef7bd5415a7ac448318e64f11a24cd21e594b/reference.go#L107-L114
		resp.JSON.Schema = "file:///" + schemaURL
	} else {
		resp.JSON.Schema = "file://" + schemaURL
	}
	return nil
}
