package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	Main(os.Args, os.Stdout, os.Stderr, NewLocalFiles("/"))
}

// The programs usage help.
var Usage = `dwi.

Usage:
	dwi <file>
`

// FileIndex provides methods for manipulating a list of files.
type FileIndex interface {
}

// LocalFiles is a FileIndex for files in the local filesystem.
type LocalFiles struct {
	Root string
}

// NewLocalFiles creates and initializes a new LocalFiles with a given path root.
func NewLocalFiles(root string) LocalFiles {
	return LocalFiles{root}
}

// Main function.
func Main(args []string, out io.Writer, eout io.Writer, fm FileIndex) {
	if len(args) < 2 {
		fmt.Fprintln(out, Usage)
		return
	}

	files := ParseArgs(args)
	fmt.Fprintln(out, files)
}

// ParseArgs takes a list of args and returns the list of files.
// The first arg must be the executable name.
func ParseArgs(args []string) []string {
	return args[1:]
}
