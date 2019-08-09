package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBooksAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books API Suite")
}
