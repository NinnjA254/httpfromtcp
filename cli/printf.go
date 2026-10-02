package cli

import (
	"flag"
	"fmt"
)

var verbose = flag.Bool("v", false, "run the program in verbose mode")

// conditional printf
func Printf(format string, a ...any) {
	if *verbose {
		fmt.Printf(format, a...)
	}
}
