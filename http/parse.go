package http

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"

	"github.com/jaypipes/gdt"
)

type jsonAssertion struct {
	Length *uint `json:"length"`
}

type responseAssertion struct {
	JSON    *jsonAssertion `json:"json"`
	Headers []string       `json:"headers"`
	Strings []string       `json:"strings"`
	Status  *int           `json:"status"`
}

type testSpec struct {
	// Name for the individual HTTP call test
	Name string `json:"name"`
	// Description of the test (defaults to Name)
	Description string `json:"description"`
	// URL being called by HTTP client
	URL string `json:"url"`
	// HTTP Method specified by HTTP client
	Method string `json:"method"`
	// Shortcut for URL and Method of "GET"
	GET string `json:"GET"`
	// Shortcut for URL and Method of "POST"
	POST string `json:"POST"`
	// JSON payload to send along in request
	Data interface{} `json:"data"`
	// Specification for expected response
	Response *responseAssertion `json:"response"`
}

type httpTestcaseConfigSchema struct {
	baseURL string
}

type httpTestcaseSchema struct {
	Config *httpTestcaseConfigSchema `json:"http"`
	Specs  []*testSpec               `json:"tests"`
}

type httpParser struct{}

// Parse accepts a Testcase and a string of YAML contents from a gdt test file.
// It then parses the HTTP test case and adds the HTTP-specific tests to the
// supplied Testcase
func (p *httpParser) Parse(tf *gdt.TestFile, contents []byte) error {
	var err error
	tcs := httpTestcaseSchema{}
	if err = yaml.Unmarshal(contents, &tcs); err != nil {
		return err
	}
	htf := &httpTestFile{
		tf, nil, nil,
	}
	for _, tspec := range tcs.Specs {
		ht := httpTest{
			tf:                htf,
			name:              tspec.Name,
			responseAssertion: tspec.Response,
		}
		if tspec.Data != nil {
			ht.jsonBody, err = json.Marshal(tspec.Data)
			if err != nil {
				return err
			}
		}
		if tspec.URL == "" {
			if tspec.GET != "" {
				ht.url = tspec.GET
				ht.method = "GET"
			} else if tspec.POST != "" {
				ht.url = tspec.POST
				ht.method = "POST"
			} else {
				return fmt.Errorf("Either specify a URL, GET or POST attribute")
			}
		} else {
			method := tspec.Method
			if method == "" {
				return fmt.Errorf("When specifying url in HTTP test spec, please specify an HTTP method")
			}
			ht.url = tspec.URL
			ht.method = method
		}
		tf.Append(&ht)
	}
	return nil
}
