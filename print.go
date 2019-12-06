package gdt

import (
	"flag"
	"fmt"
)

var optVerbosity int

func init() {
	flag.IntVar(&optVerbosity, "gdt.v", 0, "Increase verbosity of gdt library's debugging output.")
}

// V1 prints the supplied message with args if the gdt.verbosity flag has been
// set to 1 or more.
func V1(callID string, msg string, args ...interface{}) {
	debugf(1, callID, msg, args...)
}

// V2 prints the supplied message with args if the gdt.verbosity flag has been
// set to 2 or more.
func V2(callID string, msg string, args ...interface{}) {
	debugf(2, callID, msg, args...)
}

func debugf(v int, callID string, msg string, args ...interface{}) {
	if optVerbosity < v {
		return
	}
	line := fmt.Sprintf("[gdt.%s] %s", callID, msg)
	fmt.Printf(line, args...)
}
