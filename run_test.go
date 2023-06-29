// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt_test

import (
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunExecSuite(t *testing.T) {
	assert := assert.New(t)

	s, err := gdt.From("testdata/exec")
	assert.Nil(err)
	assert.NotNil(s)

	err = s.Run(gdt.NewContext(), t)
	assert.Nil(err)
}

func TestRunExecScenario(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	s, err := gdt.From("testdata/exec/ls.yaml")
	require.Nil(err)
	require.NotNil(s)

	err = s.Run(gdt.NewContext(), t)
	assert.Nil(err)
}
