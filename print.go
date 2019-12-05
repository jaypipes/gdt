package gdt

import (
	"flag"
	"fmt"
)

var optDebug bool

func init() {
	flag.BoolVar(&optDebug, "gdt.debug", false, "Turn on the gdt library's debug output.")
}

// Debugf prints the supplied message with args if the gdt.debug flag has been
// set to a truthy value
func Debugf(msg string, args ...interface{}) {
	if !optDebug {
		return
	}
	fmt.Printf(msg, args...)
}
