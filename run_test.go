// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package gdt_test

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/jaypipes/gdt-core/fixture"
	gdttypes "github.com/jaypipes/gdt-core/types"
	gdthttp "github.com/jaypipes/gdt-http"
	"github.com/jaypipes/gdt-http/test/server"
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

const (
	dataFilePath = "testdata/http/fixtures.json"
)

type dataset struct {
	Authors    interface{}
	Publishers interface{}
	Books      []*server.Book
}

func data() *dataset {
	f, err := os.Open(dataFilePath)
	if err != nil {
		panic(err)
	}
	data := &dataset{}
	if err = json.NewDecoder(f).Decode(&data); err != nil {
		panic(err)
	}
	return data
}

func dataFixture() gdttypes.Fixture {
	f, err := os.Open(dataFilePath)
	if err != nil {
		panic(err)
	}
	f.Seek(0, io.SeekStart)
	fix, err := fixture.JSON(f)
	if err != nil {
		panic(err)
	}
	return fix
}

func setup(ctx context.Context) context.Context {
	// Register an HTTP server fixture that spins up the API service on a
	// random port on localhost
	logger := log.New(os.Stdout, "books_api_http: ", log.LstdFlags)
	srv := server.NewControllerWithBooks(logger, data().Books)
	serverFixture := gdthttp.NewServerFixture(srv.Router(), false /* useTLS */)
	ctx = gdt.RegisterFixture(ctx, "books_api", serverFixture)
	ctx = gdt.RegisterFixture(ctx, "books_data", dataFixture())
	return ctx
}

func TestRunHTTPSuite(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	s, err := gdt.From("testdata/http")
	require.Nil(err)
	require.NotNil(s)

	ctx := gdt.NewContext()
	ctx = setup(ctx)

	err = s.Run(ctx, t)
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

func TestRunHTTPScenario(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	s, err := gdt.From("testdata/http/create-then-get.yaml")
	require.Nil(err)
	require.NotNil(s)

	ctx := gdt.NewContext()
	ctx = setup(ctx)

	err = s.Run(ctx, t)
	assert.Nil(err)
}
