package gdt

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

var (
	parsers = map[string]Parser{}
)

// RegisterParser registers a parser for one or more test case type names
func RegisterParser(
	parser Parser,
	types ...string,
) {
	for _, typ := range types {
		parsers[strings.ToLower(typ)] = parser
	}
}

type fileSchema struct {
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Require     []string `json:"require"`
}

// Parse reads a supplied file and parses it into a GDT File
//
// We do a double-parse of the test file. The first pass determines the
// type of test by simply looking for a "type" top-level element in the
// YAML. If no "type" element was found, the test type defaults to HTTP.
// Once the type is determined, then the test case module (e.g. gdt/http)
// is called to parse the file into the case type-specific schema
func Parse(ctx *Context, path string) (Runnable, error) {
	f, err := os.Open(path)

	if err != nil {
		panic(err)
	}
	defer f.Close()
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	tfs := fileSchema{}
	if err = yaml.Unmarshal(contents, &tfs); err != nil {
		return nil, err
	}
	parser, found := parsers[strings.ToLower(tfs.Type)]
	if !found {
		return nil, ErrUnknownParser
	}

	tf := &file{
		ctx:         ctx,
		typ:         strings.ToLower(tfs.Type),
		path:        filepath.Base(path),
		name:        tfs.Name,
		description: tfs.Description,
	}

	err = parser.Parse(tf, contents)
	if err != nil {
		return nil, err
	}

	if len(tfs.Require) > 0 {
		tf.before = make(map[string][]string, len(tfs.Require))
		// TODO(jaypipes): Parse a function-call interface from string...
		for _, elem := range tfs.Require {
			tf.before[elem] = []string{}
		}
	}

	return tf, nil
}
