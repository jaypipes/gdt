package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jaypipes/gdt/examples/books/api"
)

var _ = Describe("Books API - POST /books -> GET /books from Location", func() {

	var err error
	var resp *http.Response
	var locURL string
	var authorID, publisherID string

	Describe("proper HTTP GET after POST", func() {

		Context("when creating a single book resource", func() {
			It("should be retrievable via GET {location header}", func() {
				// See https://github.com/onsi/ginkgo/issues/457 for why this
				// needs to be here instead of in the outer Describe block.
				authorID = getAuthorByName("Ernest Hemingway").ID
				publisherID = getPublisherByName("Charles Scribner's Sons").ID
				req := api.CreateBookRequest{
					Title:       "For Whom The Bell Tolls",
					AuthorID:    authorID,
					PublisherID: publisherID,
					PublishedOn: "1940-10-21",
					Pages:       480,
				}
				var payload []byte
				payload, err = json.Marshal(&req)
				if err != nil {
					Fail("Failed to serialize JSON in setup")
				}
				resp, err = http.Post(apiPath("/books"), "application/json", bytes.NewBuffer(payload))
				Ω(err).Should(BeNil())

				// See https://github.com/onsi/ginkgo/issues/70 for why this
				// has to be one giant It() block. The GET tests rely on the
				// result of an earlier POST response (for the Location header)
				// and therefore all of the assertions below much be in a
				// single It() block. :(

				Ω(resp.StatusCode).Should(Equal(201))
				Ω(resp.Header).Should(HaveKey("Location"))

				locURL = resp.Header["Location"][0]

				resp, err = http.Get(apiPath(locURL))
				Ω(err).Should(BeNil())

				Ω(resp.StatusCode).Should(Equal(200))

				var book api.Book

				err := json.Unmarshal([]byte(respJSON(resp)), &book)
				Ω(err).Should(BeNil())

				Ω(IsValidUUID4(book.ID)).Should(BeTrue())
				Ω(book.Author).ShouldNot(BeNil())
				Ω(book.Author.Name).Should(Equal("Ernest Hemingway"))
				Ω(book.Publisher).ShouldNot(BeNil())
				Ω(book.Publisher.Address).ShouldNot(BeNil())
				Ω(book.Publisher.Address.State).Should(Equal("NY"))
			})
		})
	})
})
