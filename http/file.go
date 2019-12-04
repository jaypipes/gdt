package http

import (
	"bytes"
	"encoding/json"
	"errors"
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

var (
	errExpectedLocationHeader = errors.New("Expected Location HTTP Header in previous response")
)

type httpFileConfig struct {
	baseURL string
}

// httpFile contains groups of tests of HTTP APIs
type httpFile struct {
	ctx *gdt.Context
	cfg *httpFileConfig
	// cache of last HTTP response one of the test units executed
	PrevResponse *nethttp.Response
}

// baseURL returns the base URL to use when constructing HTTP requests
func (hf *httpFile) baseURL() string {
	// If the httpFile has been manually configured and the configuration
	// contains a base URL, use that. Otherwise, check to see if there is a
	// fixture in the registry that has an "http.base_url" state key and use
	// that if found.
	if hf.cfg != nil && hf.cfg.baseURL != "" {
		return hf.cfg.baseURL
	}
	// query the fixture registry to determine if any of them contain an
	// http.base_url state attribute.
	for _, f := range hf.ctx.Fixtures.List() {
		if f.HasState(FIXTURE_STATE_KEY_BASE_URL) {
			return f.State(FIXTURE_STATE_KEY_BASE_URL).(string)
		}
	}
	return ""
}

// client returns the HTTP client to use when executing HTTP requests. If any
// fixture provides a state with key "http.client", the fixture is asked for
// the HTTP client. Otherwise, we use the net/http.DefaultClient
func (hf *httpFile) client() *nethttp.Client {
	// query the fixture registry to determine if any of them contain an
	// http.client state attribute.
	for _, f := range hf.ctx.Fixtures.List() {
		if f.HasState(FIXTURE_STATE_KEY_CLIENT) {
			return f.State(FIXTURE_STATE_KEY_CLIENT).(*nethttp.Client)
		}
	}
	return nethttp.DefaultClient
}

// processRequestDataMap processes a map pointed to by v, transforming any
// string keys or values of the map into the results of calling the fixture
// set's State() method.
func (hf *httpFile) preprocessMap(
	m reflect.Value, kt reflect.Type, vt reflect.Type,
) error {
	it := m.MapRange()
	for it.Next() {
		if kt.Kind() == reflect.String {
			keyStr := it.Key().String()
			for _, f := range hf.ctx.Fixtures.List() {
				if !f.HasState(keyStr) {
					continue
				}
				trKeyStr := f.State(keyStr)
				keyStr = trKeyStr.(string)
			}

			val := it.Value()
			err := hf.preprocessMapValue(m, reflect.ValueOf(keyStr), val, val.Type())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (hf *httpFile) preprocessMapValue(m reflect.Value, k reflect.Value, v reflect.Value, vt reflect.Type) error {
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
		return hf.preprocessMap(v, vt.Key(), vt.Elem())
	case reflect.String:
		valStr := v.String()
		for _, f := range hf.ctx.Fixtures.List() {
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

// httpTest represents a single HTTP request and response pair along with
// expectations/assertions for the response components
type httpTest struct {
	f *httpFile
	// Name for the individual HTTP call test
	name string
	// Description of the test (defaults to Name)
	description string
	// URL being called by HTTP client
	url string
	// HTTP Method specified by HTTP client
	method string
	// Data to send in the request, typically serialized as JSON. This data is
	// pre-processed to replace values that look like JSONPath expressions with
	// fixture data.
	data interface{}
	// Specification for expected response
	responseAssertion *responseAssertion
}

// getURL returns the URL to use for the test's HTTP request. The test's url
// field is first queried to see if it is the special $LOCATION string. If it
// is, then we return the previous HTTP response's Location header. Otherwise,
// we construct the URL from the httpFile's base URL and the test's url field.
func (ht *httpTest) getURL() (string, error) {
	if strings.ToUpper(ht.url) == "$LOCATION" {
		if ht.f.PrevResponse == nil {
			panic("test unit referenced $LOCATION before executing an HTTP request")
		}
		url, err := ht.f.PrevResponse.Location()
		if err != nil {
			return "", errExpectedLocationHeader
		}
		return url.String(), nil
	}
	base := ht.f.baseURL()
	return base + ht.url, nil
}

// processRequestData looks through the raw data interface{} that was
// unmarshaled during parse for any string values that look like JSONPath
// expressions. If we find any, we query the fixture registry to see if any
// fixtures have a value that matches the JSONPath expression. See
// gdt.fixtures:jsonFixture for more information on how this works
func (ht *httpTest) processRequestData() {
	if ht.data == nil {
		return
	}

	// Get a pointer to the unmarshaled interface{} so we can mutate the
	// contents pointed to
	p := reflect.ValueOf(&ht.data)

	// We're interested in the value pointed to by the interface{}, which is
	// why we do a double Elem() here.
	v := p.Elem().Elem()
	vt := v.Type()

	switch vt.Kind() {
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Elem()
			it := item.Type()
			ht.f.preprocessMap(item, it.Key(), it.Elem())
		}
		//	ht.f.preprocessSliceValue(v, vt.Key(), vt.Elem())
	case reflect.Map:
		ht.f.preprocessMap(v, vt.Key(), vt.Elem())
	}
}

// Run executes the test described by the HTTP test. A new HTTP request and
// response pair is created during this call.
func (ht *httpTest) Run(t *testing.T) {
	var body io.Reader
	if ht.data != nil {
		ht.processRequestData()
		jsonBody, err := json.Marshal(ht.data)
		require.Nil(t, err)
		body = bytes.NewReader(jsonBody)
	}
	t.Run(ht.name, func(t *testing.T) {
		url, err := ht.getURL()
		require.Nil(t, err)
		req, err := nethttp.NewRequest(ht.method, url, body)
		require.Nil(t, err)
		// TODO(jaypipes): Allow customization of the HTTP client for proxying,
		// TLS, etc
		c := ht.f.client()
		resp, err := c.Do(req)
		require.Nil(t, err)
		if ht.responseAssertion != nil {
			// Only read the response body contents once and pass the byte
			// buffer to the assertion functions
			b, err := ioutil.ReadAll(resp.Body)
			require.Nil(t, err)

			rspec := ht.responseAssertion
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
		ht.f.PrevResponse = resp
	})
}
