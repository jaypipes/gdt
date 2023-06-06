// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	nethttp "net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/jaypipes/gdt"
	"github.com/stretchr/testify/require"
)

// TestSpec describes a a test of a single HTTP request and response
type TestSpec struct {
	defaults *TestCaseDefaults
	// Name for the individual HTTP call test
	Name string `json:"name,omitempty"`
	// Description of the test (defaults to Name)
	Description string `json:"description,omitempty"`
	// URL being called by HTTP client
	URL string `json:"url,omitempty"`
	// HTTP Method specified by HTTP client
	Method string `json:"method,omitempty"`
	// Shortcut for URL and Method of "GET"
	GET string `json:"GET,omitempty"`
	// Shortcut for URL and Method of "POST"
	POST string `json:"POST,omitempty"`
	// Shortcut for URL and Method of "PUT"
	PUT string `json:"PUT,omitempty"`
	// Shortcut for URL and Method of "PATCH"
	PATCH string `json:"PATCH,omitempty"`
	// Shortcut for URL and Method of "DELETE"
	DELETE string `json:"DELETE,omitempty"`
	// JSON payload to send along in request
	Data interface{} `json:"data,omitempty"`
	// Specification for expected response
	Response *ResponseAssertion `json:"response,omitempty"`
}

// getURL returns the URL to use for the test's HTTP request. The test's url
// field is first queried to see if it is the special $LOCATION string. If it
// is, then we return the previous HTTP response's Location header. Otherwise,
// we construct the URL from the httpFile's base URL and the test's url field.
func (s *TestSpec) getURL(ctx context.Context) (string, error) {
	if strings.ToUpper(s.URL) == "$LOCATION" {
		pr := getPreviousResponse(ctx)
		if pr == nil {
			panic("test unit referenced $LOCATION before executing an HTTP request")
		}
		url, err := pr.Location()
		if err != nil {
			return "", ErrExpectedLocationHeader
		}
		return url.String(), nil
	}
	base := s.defaults.BaseURLFromContext(ctx)
	return base + s.URL, nil
}

// processRequestData looks through the raw data interface{} that was
// unmarshaled during parse for any string values that look like JSONPath
// expressions. If we find any, we query the fixture registry to see if any
// fixtures have a value that matches the JSONPath expression. See
// gdt.fixtures:jsonFixture for more information on how this works
func (s *TestSpec) processRequestData(ctx context.Context) {
	if s.Data == nil {
		return
	}
	gdt.V3("http.file.TestSpec:processRequestData", "ht.data: %+v\n", s.Data)
	// Get a pointer to the unmarshaled interface{} so we can mutate the
	// contents pointed to
	p := reflect.ValueOf(&s.Data)

	// We're interested in the value pointed to by the interface{}, which is
	// why we do a double Elem() here.
	v := p.Elem().Elem()
	vt := v.Type()

	switch vt.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Elem()
			it := item.Type()
			s.preprocessMap(ctx, item, it.Key(), it.Elem())
		}
		//	ht.f.preprocessSliceValue(v, vt.Key(), vt.Elem())
	case reflect.Map:
		s.preprocessMap(ctx, v, vt.Key(), vt.Elem())
	}
}

// client returns the HTTP client to use when executing HTTP requests. If any
// fixture provides a state with key "http.client", the fixture is asked for
// the HTTP client. Otherwise, we use the net/http.DefaultClient
func (s *TestSpec) client(ctx context.Context) *nethttp.Client {
	// query the fixture registry to determine if any of them contain an
	// http.client state attribute.
	for _, f := range gdt.GetFixturesFromContext(ctx).List() {
		if f.HasState(StateKeyClient) {
			c, ok := f.State(StateKeyClient).(*nethttp.Client)
			if !ok {
				panic("fixture failed to return a *net/http.Client")
			}
			return c
		}
	}
	return nethttp.DefaultClient
}

// processRequestDataMap processes a map pointed to by v, transforming any
// string keys or values of the map into the results of calling the fixture
// set's State() method.
func (s *TestSpec) preprocessMap(
	ctx context.Context,
	m reflect.Value,
	kt reflect.Type,
	vt reflect.Type,
) error {
	it := m.MapRange()
	for it.Next() {
		if kt.Kind() == reflect.String {
			keyStr := it.Key().String()
			for _, f := range gdt.GetFixturesFromContext(ctx).List() {
				if !f.HasState(keyStr) {
					continue
				}
				trKeyStr := f.State(keyStr)
				keyStr = trKeyStr.(string)
			}

			val := it.Value()
			err := s.preprocessMapValue(ctx, m, reflect.ValueOf(keyStr), val, val.Type())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *TestSpec) preprocessMapValue(
	ctx context.Context,
	m reflect.Value,
	k reflect.Value,
	v reflect.Value,
	vt reflect.Type,
) error {
	if vt.Kind() == reflect.Interface {
		v = v.Elem()
		vt = v.Type()
	}

	switch vt.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			fmt.Println(item)
		}
		fmt.Printf("map element is an array.\n")
	case reflect.Map:
		return s.preprocessMap(ctx, v, vt.Key(), vt.Elem())
	case reflect.String:
		valStr := v.String()
		for _, f := range gdt.GetFixturesFromContext(ctx).List() {
			if !f.HasState(valStr) {
				continue
			}
			trValStr := f.State(valStr)
			m.SetMapIndex(k, reflect.ValueOf(trValStr))
		}
	default:
		return nil
	}
	return nil
}

// Run executes the test described by the HTTP test. A new HTTP request and
// response pair is created during this call.
func (s *TestSpec) Run(ctx context.Context, t *testing.T) context.Context {
	var body io.Reader
	if s.Data != nil {
		s.processRequestData(ctx)
		jsonBody, err := json.Marshal(s.Data)
		require.Nil(t, err)
		body = bytes.NewReader(jsonBody)
	}
	t.Run(s.Name, func(t *testing.T) {
		url, err := s.getURL(ctx)
		if err != nil {
			panic(err)
		}

		req, err := nethttp.NewRequest(s.Method, url, body)
		if err != nil {
			panic(err)
		}

		// TODO(jaypipes): Allow customization of the HTTP client for proxying,
		// TLS, etc
		c := s.client(ctx)

		resp, err := c.Do(req)
		if err != nil {
			panic(err)
		}

		// Make sure we drain and close our response body...
		defer func() {
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()

		if s.Response != nil {
			// Only read the response body contents once and pass the byte
			// buffer to the assertion functions
			b, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)

			rspec := s.Response
			if rspec.Status != nil {
				assertHTTPStatusEqual(t, resp, *(rspec.Status))
			}

			if rspec.JSON != nil {
				assertJSON(t, resp, b, rspec.JSON)
			}

			if len(rspec.Strings) > 0 {
				for _, exp := range rspec.Strings {
					assertStringInBody(t, resp, b, exp)
				}
			}

			if len(rspec.Headers) > 0 {
				for _, exp := range rspec.Headers {
					assertHeader(t, resp, exp)
				}
			}
		}
		ctx = storePreviousResponse(ctx, resp)
	})
	return ctx
}
