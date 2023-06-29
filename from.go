// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"io"
	"os"

	gdterrors "github.com/jaypipes/gdt-core/errors"
	"github.com/jaypipes/gdt-core/scenario"
	"github.com/jaypipes/gdt-core/suite"
	gdttypes "github.com/jaypipes/gdt-core/types"
)

// From returns a new Runnable thing from an `io.Reader`, a string file or
// directory path, or the raw bytes of YAML content describing a scenario or
// suite.
func From(source interface{}) (gdttypes.Runnable, error) {
	defaultContext := NewContext()
	switch source.(type) {
	case io.Reader:
		return scenario.FromReader(
			source.(io.Reader),
			scenario.WithContext(defaultContext),
		)
	case string:
		path := source.(string)
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if fi.IsDir() {
			return suite.FromDir(
				path,
				suite.WithContext(defaultContext),
			)
		} else {
			return scenario.FromReader(
				f,
				scenario.WithPath(path),
				scenario.WithContext(defaultContext),
			)
		}
	case []byte:
		return scenario.FromBytes(source.([]byte))
	default:
		return nil, gdterrors.UnknownSourceType(source)
	}
}
