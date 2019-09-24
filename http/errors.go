package http

import "errors"

var (
	// ErrInvalidAliasOrURL is returned when the test author failed to provide
	// either a URL and Method or specify one of the aliases like GET, POST, or
	// DELETE
	// TODO(jaypipes): wrap with gdt.ErrInvalid
	ErrInvalidAliasOrURL = errors.New("Either specify a URL and Method or specify one of GET, POST, PUT, PATCH or DELETE")
)
