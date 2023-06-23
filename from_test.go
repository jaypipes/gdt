// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt_test

import (
	"os"
	"testing"

	"github.com/jaypipes/gdt"
	gdterrors "github.com/jaypipes/gdt-core/errors"
	"github.com/jaypipes/gdt-core/scenario"
	"github.com/jaypipes/gdt-core/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromUnknownSourceType(t *testing.T) {
	assert := assert.New(t)

	s, err := gdt.From(1)
	assert.NotNil(err)
	assert.Nil(s)

	assert.ErrorIs(err, gdterrors.ErrUnknownSourceType)
}

func TestFromFileNotFound(t *testing.T) {
	assert := assert.New(t)

	s, err := gdt.From("/path/to/nonexisting/file")
	assert.NotNil(err)
	assert.Nil(s)

	assert.True(os.IsNotExist(err))
}

func TestFromSuite(t *testing.T) {
	assert := assert.New(t)

	s, err := gdt.From("testdata/exec")
	assert.Nil(err)
	assert.NotNil(s)

	suite, ok := s.(*suite.Suite)
	assert.True(ok, "gdt.From() did not return a Suite")

	assert.Equal("testdata/exec", suite.Path)
	assert.Len(suite.Scenarios, 2)
}

func TestFromScenarioPath(t *testing.T) {
	assert := assert.New(t)

	s, err := gdt.From("testdata/exec/ls.yaml")
	assert.Nil(err)
	assert.NotNil(s)

	sc, ok := s.(*scenario.Scenario)
	assert.True(ok, "gdt.From() with dir path did not return a Scenario")

	assert.Equal("testdata/exec/ls.yaml", sc.Path)
	assert.Len(sc.Tests, 1)
}

func TestFromScenarioReader(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, err := os.Open("testdata/exec/ls.yaml")
	require.Nil(err)
	s, err := gdt.From(f)
	assert.Nil(err)
	assert.NotNil(s)

	sc, ok := s.(*scenario.Scenario)
	assert.True(ok, "gdt.From() from file path did not return a Scenario")

	// The scenario's path isn't set because we didn't supply a filepath...
	assert.Equal("", sc.Path)
	assert.Len(sc.Tests, 1)
}

func TestFromScenarioBytes(t *testing.T) {
	assert := assert.New(t)

	raw := `name: foo
description: simple foo test
tests:
 - exec: echo foo
`
	b := []byte(raw)
	s, err := gdt.From(b)
	assert.Nil(err)
	assert.NotNil(s)

	sc, ok := s.(*scenario.Scenario)
	assert.True(ok, "gdt.From() with []byte did not return a Scenario")

	// The scenario's path isn't set because we didn't supply a filepath...
	assert.Equal("", sc.Path)
	assert.Len(sc.Tests, 1)
}
