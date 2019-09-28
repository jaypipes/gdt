package http

import (
	"github.com/ghodss/yaml"

	"github.com/jaypipes/gdt"
)

type jsonAssertion struct {
	Length      *uint             `json:"length"`
	Paths       map[string]string `json:"paths"`
	PathFormats map[string]string `json:"path_formats"`
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
	// Shortcut for URL and Method of "PUT"
	PUT string `json:"PUT"`
	// Shortcut for URL and Method of "PATCH"
	PATCH string `json:"PATCH"`
	// Shortcut for URL and Method of "DELETE"
	DELETE string `json:"DELETE"`
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
func (p *httpParser) Parse(ca gdt.ContextAppendable, contents []byte) error {
	var err error
	tcs := httpTestcaseSchema{}
	if err = yaml.Unmarshal(contents, &tcs); err != nil {
		return err
	}
	hf := &httpFile{
		ca.Context(), nil, nil,
	}
	for _, tspec := range tcs.Specs {
		ht := httpTest{
			f:                 hf,
			name:              tspec.Name,
			responseAssertion: tspec.Response,
			data:              tspec.Data,
		}
		ht.method, ht.url, err = parseMethodAndURL(tspec)
		if err != nil {
			return err
		}
		ca.Append(&ht)
	}
	return nil
}

func parseMethodAndURL(tspec *testSpec) (string, string, error) {
	if tspec.URL == "" {
		if tspec.GET != "" {
			return "GET", tspec.GET, nil
		} else if tspec.POST != "" {
			return "POST", tspec.POST, nil
		} else if tspec.PUT != "" {
			return "PUT", tspec.PUT, nil
		} else if tspec.DELETE != "" {
			return "DELETE", tspec.DELETE, nil
		} else {
			return "", "", ErrInvalidAliasOrURL
		}
	}
	if tspec.Method == "" {
		return "", "", ErrInvalidAliasOrURL
	}
	return tspec.Method, tspec.URL, nil
}
