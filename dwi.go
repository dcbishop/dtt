package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/cep21/xdgbasedir"
	"gopkg.in/yaml.v1"
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
	Open(filename string) (*os.File, error)
}

// LocalFiles is a FileIndex for files in the local filesystem.
type LocalFiles struct {
	Root string
}

// Exists returns true if a file given by path exists.
func (lf *LocalFiles) Exists(filename string) bool {
	if _, err := os.Stat(lf.getFullPath(filename)); err != nil {
		return false
	}
	return true
}

// Open opens a file or returns an error.
func (lf *LocalFiles) Open(filename string) (*os.File, error) {
	return os.Open(filename)
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

	configDir, err := xdgbasedir.ConfigHomeDirectory()
	if err != nil {
		return
	}

	rules := LoadRules(path.Join(configDir, "/dealwithit/rules.yaml"), fi, eout)

	log.Println(rules)
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

// LoadRules loads the rules file given by filename from the given FileIndex.
func LoadRules(filename string, fi FileIndex, eout io.Writer) []map[string]string {
	var rules []map[string]string
	file, err := fi.Open(filename)
	if err != nil {
		fmt.Fprintln(eout, "Error: Could not open rules file", filename, err)
		return rules
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintln(eout, "Error: Could not read rules file", filename, err)
		return rules
	}

	err = yaml.Unmarshal(data, &rules)
	if err != nil {
		fmt.Fprintln(eout, "Error: Could not parse rules file", filename, err)
		return rules
	}

	return rules
}
