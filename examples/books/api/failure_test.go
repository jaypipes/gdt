package api_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	defaultAPIServerURL = "http://localhost:8081"
)

// respJSON returns a string if the supplied HTTP response body is JSON,
// otherwise the empty string
func respJSON(r *http.Response) string {
	if r == nil {
		return ""
	}
	if !strings.HasPrefix(r.Header.Get("content-type"), "application/json") {
		return ""
	}
	bodyStr, _ := ioutil.ReadAll(r.Body)
	return string(bodyStr)
}

// respText returns a string if the supplied HTTP response has a text/plain
// content type and a body, otherwise the empty string
func respText(r *http.Response) string {
	if r == nil {
		return ""
	}
	if !strings.HasPrefix(r.Header.Get("content-type"), "text/plain") {
		return ""
	}
	bodyStr, _ := ioutil.ReadAll(r.Body)
	return string(bodyStr)
}

func apiPath(path string) string {
	serverURL, found := os.LookupEnv("EXAMPLES_BOOKS_API_SERVER_URL")
	if !found {
		serverURL = defaultAPIServerURL
	}
	return strings.TrimSuffix(serverURL, "/") + "/" + strings.TrimPrefix(path, "/")
}

var _ = Describe("Books API - GET /books failures", func() {
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

		Context("when no such book was found", func() {
			JustBeforeEach(func() {
				testPath = "/books/nosuchbook"
			})

			AssertZeroJSONLength()

			It("should return 404", func() {
				Ω(response.StatusCode).Should(Equal(404))
			})
		})

		Context("when an invalid query parameter is supplied", func() {
			JustBeforeEach(func() {
				testPath = "/books?invalidparam=1"
			})

			AssertZeroJSONLength()

			It("should return 400", func() {
				Ω(response.StatusCode).Should(Equal(400))
			})
			It("should indicate invalid query parameter", func() {
				Ω(respText(response)).Should(ContainSubstring("invalid parameter"))
			})
		})
	})
})
