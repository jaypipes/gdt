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
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := &Context{
		Fixtures: Fixtures,
	}

	fi, err := f.Stat()
	switch {
	case err != nil:
		panic(err)
	case fi.IsDir():
		{
			// List YAML files in the directory and parse each into a testable unit
			var files []string
			s := &suite{
				path: path,
				// TODO(jaypipes): Allows name/description of suite
				name:        path,
				description: path,
			}

			err := filepath.Walk(path, func(subpath string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				suffix := filepath.Ext(subpath)
				if suffix != ".yaml" {
					return nil
				}
				files = append(files, subpath)
				return nil
			})
			if err != nil {
				panic(err)
			}
			for _, fp := range files {
				tf, err := Parse(ctx, fp)
				if err != nil {
					panic(err)
				}
				s.Append(tf)
			}
			return s, nil
		}
	default:
		tf, err := Parse(ctx, path)
		if err != nil {
			panic(err)
		}
		return tf, nil
	}
}
