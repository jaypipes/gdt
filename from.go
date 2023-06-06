// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"os"
	"path/filepath"
)

// From returns a Runnable thing after reading a supplied filepath and
// parsing the file or directory into a test case or test suite
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
		return NewTestSuiteFromDir(path)
	default:
		tc, err := NewTestCaseFromReader(f, path)
		if err != nil {
			return nil, err
		}
		return tc, nil
	}
}
