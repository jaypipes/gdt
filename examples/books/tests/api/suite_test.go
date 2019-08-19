package api_test

import (
	"fmt"
	"testing"

	"github.com/jaypipes/gdt"
	_ "github.com/jaypipes/gdt/http"
)

func TestBooksAPI(t *testing.T) {
	tc, err := gdt.FromFile("failures.yaml")
	if err != nil {
		panic(err)
	}
	res := tc.Run(nil, nil)
	fmt.Println(res)
}
