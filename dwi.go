package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	Main(os.Args, os.Stdout, os.Stderr)
}

// The programs usage help.
var Usage = `dwi.

Usage:
	dwi <file>
`

// Main function.
func Main(args []string, out io.Writer, eout io.Writer) {
	if len(args) < 2 {
		fmt.Fprintln(out, Usage)
		return
	}
	fmt.Fprintln(out, "Your filename was:", args[1])
}
