package http

import (
	"fmt"

	"github.com/google/uuid"
)

type validator func(string) bool

var (
	validators = map[string]validator{
		"uuid":  isUUID,
		"uuid4": isUUID4,
	}
)

// isFormatted takes a format string and a string value, determines the
// validator function for that type of format string and returns whether the
// value string is formatted correctly.
func isFormatted(format string, value string) (bool, error) {
	fn, ok := validators[format]
	if !ok {
		return false, fmt.Errorf("unknown format %s", format)
	}
	return fn(value), nil
}

func isUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func isUUID4(s string) bool {
	u, err := uuid.Parse(s)
	if err != nil {
		return false
	}
	return u.Version() == 4
}
