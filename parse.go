package gdt

import (
	"io"
	"io/ioutil"
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

// parse parses the supplied file and returns an error if any syntax or
// validation failed
//
// We do a double-parse of the test file. The first pass determines the
// type of test by simply looking for a "type" top-level element in the
// YAML. If no "type" element was found, the test type defaults to HTTP.
// Once the type is determined, then the test case module (e.g. gdt/http)
// is called to parse the file into the case type-specific schema
func (tf *file) parse(r io.Reader) error {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return parseBytes(tf, contents, &parsers)
}

func parseBytes(
	tf *file,
	contents []byte,
	typeParsers *map[string]Parser,
) error {
	tfs := fileSchema{}
	if err := yaml.Unmarshal(contents, &tfs); err != nil {
		return ErrInvalidYAML
	}
	parser, found := (*typeParsers)[strings.ToLower(tfs.Type)]
	if !found {
		return ErrUnknownParser
	}

	tf.typ = strings.ToLower(tfs.Type)
	tf.name = tfs.Name
	tf.description = tfs.Description
	tf.require = tfs.Require

	if err := parser.Parse(tf, contents); err != nil {
		return err
	}

	return nil
}
