package api_test

import (
	"github.com/jaypipes/gdt/examples/books/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Books API Types", func() {
	var (
		longBook  api.Book
		shortBook api.Book
	)

	BeforeEach(func() {
		longBook = api.Book{
			Title: "Les Miserables",
			Pages: 1488,
			Author: &api.Author{
				Name: "Victor Hugo",
			},
		}

		shortBook = api.Book{
			Title: "Fox In Socks",
			Pages: 24,
			Author: &api.Author{
				Name: "Dr. Seuss",
			},
		}
	})

	Describe("Categorizing book length", func() {
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
			})
		})

		Context("With fewer than 300 pages", func() {
			It("should be a short story", func() {
				Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
			})
		})
	})
})
