package http

import "github.com/jaypipes/gdt"

func init() {
	gdt.RegisterParser(&httpParser{}, "http", "")
}
