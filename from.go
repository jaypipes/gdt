// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt

import (
	"os"
	"path/filepath"
)

// From returns a Runnable thing after reading a supplied filepath and
// parsing the file or directory into a test file or test suite
func From(path string) (Runnable, error) {
	// Determine if the path is a directory or a regular file. If it's a
	// directory, construct a suite. If it's a regular file, construct a test
	// file by parsing the contents.
	path, _ = filepath.Abs(path)
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	ctx := &Context{
		Fixtures: Fixtures,
	}

	fi, err := f.Stat()
	switch {
	case err != nil:
		return nil, err
	case fi.IsDir():
		ctx.Basedir = path
		return fromDir(ctx, path)
	default:
		ctx.Basedir = filepath.Dir(path)

		tf := &file{
			ctx:  ctx,
			path: filepath.Base(path),
		}
		if err = tf.parse(f); err != nil {
			return nil, err
		}
		return tf, nil
	}
}

func fromDir(ctx *Context, dirPath string) (Runnable, error) {
	// List YAML files in the directory and parse each into a testable unit
	s := &suite{
		path: dirPath,
		// TODO(jaypipes): Allows name/description of suite
		name:        dirPath,
		description: dirPath,
	}

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

			tf := &file{
				ctx:  ctx,
				path: filepath.Base(path),
			}
			if err = tf.parse(f); err != nil {
				return err
			}
			s.Append(tf)
			return nil
		},
	); err != nil {
		return nil, err
	}
	return s, nil
}
