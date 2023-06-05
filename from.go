// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

// From returns a Runnable thing after reading a supplied filepath and
// parsing the file or directory into a test file or test suite
func From(path string) (Runnable, error) {
	// Determine if the path is a directory or a regular file. If it's a
	// directory, construct a suite. If it's a regular file, construct a test
	// case by parsing the contents.
	path, _ = filepath.Abs(path)
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	switch {
	case err != nil:
		return nil, err
	case fi.IsDir():
		return FromDir(path)
	default:
		tc, err := FromReader(f)
		if err != nil {
			return nil, err
		}
		return tc, nil
	}
}

// FromDir reads the supplied directory path and returns a Runnable
// representing the suite of test cases in that directory.
func FromDir(dirPath string) (Runnable, error) {
	// List YAML files in the directory and parse each into a testable unit
	s := NewTestSuite(WithPath(dirPath))

	if err := filepath.Walk(
		dirPath,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			suffix := filepath.Ext(path)
			if suffix != ".yaml" {
				return nil
			}
			f, err := os.Open(path)

			if err != nil {
				return err
			}
			defer f.Close()

			tc, err := FromReader(f)
			if err != nil {
				return err
			}
			s.Append(tc)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return s, nil
}

// FromReader parses the supplied io.Reader and returns a Runnable representing
// the contents in the reader. Returns an error if any syntax or validation
// failed
func FromReader(r io.Reader) (Runnable, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return FromBytes(contents)
}

// FromBytes returns a Runnable after processing the supplied contents
// representing a test case
func FromBytes(contents []byte) (Runnable, error) {
	// We do a double-parse of the test file. The first pass determines the
	// type of test by simply looking for a "type" top-level element in the
	// YAML. If no "type" element was found, the test type defaults to HTTP.
	// Once the type is determined, then the test case module (e.g. gdt/http)
	// is called to parse the file into the case type-specific schema
	tc := TestCase{}
	if err := yaml.Unmarshal(contents, &tc); err != nil {
		return nil, ErrInvalidYAML
	}

	tc.Type = strings.ToLower(tc.Type)
	parser := TestTypeParsers.Get(tc.Type)
	if parser == nil {
		return nil, ErrUnknownParser
	}

	if err := parser.Parse(&tc, contents); err != nil {
		return nil, err
	}

	return &tc, nil
}
