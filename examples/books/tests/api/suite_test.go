package api_test

import (
	"testing"

	"github.com/jaypipes/gdt"
	_ "github.com/jaypipes/gdt/http"
)

func TestBooksAPI(t *testing.T) {
	gdt.FromFile(t, "failures.yaml")
}
