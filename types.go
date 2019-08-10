package api_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestType uint

const (
	TestTypeHTTP TestType = iota
)

type TestAction struct {
	Name        string
	Description string
	Before      *TestAction
}

type testTypeProber struct {
	Type string `json:"type"`
}

type TestBase struct {
	// Type of test this is
	Type TestType
	// Filepath is the filepath to the test file
	Filepath string
	// Name for the overall test
	Name string `json:"name"`
	// Description of the test (defaults to Name)
	Description string `json:"description"`
}

type JSONResponseAssertion struct {
	Length uint `json:"length"`
}

type HTTPResponseAssertion struct {
	JSON    *JSONResponseAssertion `json:"json"`
	Strings []string               `json:"strings"`
	Status  int                    `json:"status"`
}

type HTTPTestSpec struct {
	// Name for the individual HTTP call test
	Name string `json:"name"`
	// Description of the test (defaults to Name)
	Description string `json:"description"`
	// URL being called by HTTP client
	URL string `json:"url"`
	// HTTP Method specified by HTTP client
	Method string `json:"method"`
	// Shortcut for URL and Method of "GET"
	GET string
	// Shortcut for URL and Method of "POST"
	POST string
	// HTTP request object constructed from spec
	Request *http.Request
	// Specification for expected response
	Response *HTTPResponseAssertion `json:"response"`
}

type HTTPTest struct {
	*TestBase
	TestSpecs []*HTTPTestSpec `json:"tests"`
}

// TestFromFile reads a GDT test from the supplied filepath and creates Ginkgo
// test elements
func TestFromFile(fp string) error {
	// We do a double-parse of the test file. The first pass determines the
	// type of test by simply looking for a "type" element in the YAML. If no
	// "type" element was found, the test type defaults to HTTP. Once the type
	// is determined, then the test file is unmarshaled into the concrete
	// $TYPETest struct.
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	tp := testTypeProber{}
	if err = yaml.Unmarshal(contents, &tp); err != nil {
		return err
	}

	switch tp.Type {
	case "http", "":
		{
			ht := HTTPTest{}
			if err := yaml.Unmarshal(contents, &ht); err != nil {
				return err
			}
			Describe("books API failures", func() {
				var response *http.Response
				var err error
				var testPath string

				BeforeEach(func() {
					response, err = http.Get(apiPath(testPath))
					Ω(err).Should(BeZero())
				})

				Describe("failure modes", func() {

					AssertZeroJSONLength := func() {
						It("should not include JSON in the response", func() {
							Ω(respJSON(response)).Should(BeZero())
						})
					}

					for _, tspec := range ht.TestSpecs {
						Context(tspec.Name, func() {
							JustBeforeEach(func() {
								testPath = tspec.GET
							})

							if tspec.Response != nil {
								rspec := tspec.Response
								if rspec.JSON != nil {
									if rspec.JSON.Length == 0 {
										AssertZeroJSONLength()
									}
								}

								if rspec.Status != 0 {
									It(fmt.Sprintf("should return %d", rspec.Status), func() {
										Ω(response.StatusCode).Should(Equal(rspec.Status))
									})
								}

								if len(rspec.Strings) > 0 {
									for _, expStr := range rspec.Strings {
										It(fmt.Sprintf("should contain '%s'", expStr), func() {
											Ω(respText(response)).Should(ContainSubstring(expStr))
										})
									}
								}
							}
						})
					}
				})
			})
			return nil
		}
	default:
		return fmt.Errorf("Unknown test type specified: %s", tp.Type)
	}
}
