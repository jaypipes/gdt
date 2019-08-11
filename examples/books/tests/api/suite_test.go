package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gdt "../../../.."
)

func TestBooksAPI(t *testing.T) {
	RegisterFailHandler(Fail)
	err := gdt.TestFromFile("failures.yaml")
	if err != nil {
		panic(err)
	}
	RunSpecs(t, "Books API Suite")
}
