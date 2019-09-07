package api_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
