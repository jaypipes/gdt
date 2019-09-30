package http

import (
	"github.com/pkg/errors"

	"github.com/jaypipes/gdt"
)

var (
	// ErrInvalidAliasOrURL is returned when the test author failed to provide
	// either a URL and Method or specify one of the aliases like GET, POST, or
	// DELETE
	ErrInvalidAliasOrURL = errors.Wrap(
		gdt.ErrInvalid,
		"Either specify a URL and Method or specify one of GET, POST, PUT, PATCH or DELETE",
	)
)
