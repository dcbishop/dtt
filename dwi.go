package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func main() {
	lf := NewLocalFiles("")
	Main(os.Args, os.Stdout, os.Stderr, &lf)
}

// The programs usage help.
var Usage = `dwi.

Usage:
	dwi <file>
`

// FileIndex provides methods for manipulating a list of files.
type FileIndex interface {
	Exists(filename string) bool
}

// LocalFiles is a FileIndex for files in the local filesystem.
type LocalFiles struct {
	Root string
}

// Exists returns true if a file given by path exists.
func (lf *LocalFiles) Exists(filename string) bool {
	if _, err := os.Stat(lf.getFullPath(filename)); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetFullPath returns the full path.
func (lf *LocalFiles) getFullPath(filename string) string {
	return path.Join(lf.Root, filename)
}

// NewLocalFiles creates and initializes a new LocalFiles with a given path root.
func NewLocalFiles(root string) LocalFiles {
	return LocalFiles{root}
}

// Main function.
func Main(args []string, out io.Writer, eout io.Writer, fi FileIndex) {
	if len(args) < 2 {
		fmt.Fprintln(out, Usage)
		return
	}

	files := ParseArgs(args)
	if AnyFilesDontExist(eout, fi, files) {
		return
	}
}

// AnyFilesDontExist will look for each of the given files in the given FileIndex,
// output an error to eout for any file that doesn't exist and return true if where any missing.
func AnyFilesDontExist(eout io.Writer, fi FileIndex, files []string) bool {
	err := false
	for _, f := range files {
		if !fi.Exists(f) {
			fmt.Fprintln(eout, "Error: File not found:", f)
			err = true
		}
	}
	return err
}

// ParseArgs takes a list of args and returns the list of files.
// The first arg must be the executable name.
func ParseArgs(args []string) []string {
	return args[1:]
}
