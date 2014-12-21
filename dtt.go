package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/cep21/xdgbasedir"
	"github.com/termie/go-shutil"
	"gopkg.in/yaml.v1"
)

func main() {
	lf := NewLocalFiles("")
	Main(os.Args, os.Stdout, os.Stderr, &lf)
}

// The programs usage help.
var Usage = `dtt.

Usage:
	dtt <file>
`

// FileIndex provides methods for manipulating a list of files.
type FileIndex interface {
	Exists(filename string) bool
	Open(filename string) (*os.File, error)
	IsDir(filename string) bool
	Rename(filename string, dest string) error
	Move(filename string, dest string) error
	Remove(filename string) error
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

// IsDir returns true if filename points to a directory.
func (lf *LocalFiles) IsDir(filename string) bool {
	fileInfo, err := os.Stat(lf.getFullPath(filename))
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// Open opens a file or returns an error.
func (lf *LocalFiles) Open(filename string) (*os.File, error) {
	return os.Open(lf.getFullPath(filename))
}

// Rename renames (moves) a file.
func (lf *LocalFiles) Rename(oldpath, newpath string) error {
	return os.Rename(lf.getFullPath(oldpath), lf.getFullPath(newpath))

}

// Move moves a directory or file to the destination directory.
func (lf *LocalFiles) Move(oldpath, newpath string) error {
	cmd := exec.Command("mv", "-v", lf.getFullPath(oldpath), lf.getFullPath(newpath))
	out, err := cmd.Output()
	fmt.Println(string(out))
	return err

	//err := lf.Rename(oldpath, newpath)

	//switch err.(type) {
	//case *os.LinkError:
	//err = lf.CopyWithTemp(oldpath, newpath)
	//}

	//return err
}

// Remove removes a file.
func (lf *LocalFiles) Remove(filename string) error {
	return os.Remove(lf.getFullPath(filename))
}

// CopyWithTemp copies a directory to the destination.
func (lf *LocalFiles) CopyWithTemp(oldpath, newpath string) error {
	// [TODO]: Actually do temp stuff - 2014-12-21 01:49pm
	src := lf.getFullPath(oldpath)
	dst := lf.getFullPath(newpath)
	//dstTemp := dst + ".diw-COPYING"
	//err := shutil.CopyTree(src, dstTemp, nil)
	err := shutil.CopyTree(src, dst, nil)
	if err != nil {
		return err
	}
	return err

	//return os.Rename(dstTemp, dst)
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

	rules := LoadRules(path.Join(configDir, "/dothething/rules.yaml"), fi, eout)

	ForEachMatchingFile(files, rules, out, eout, fi)
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

type Rule map[string]string
type Rules []Rule

// LoadRules loads the rules file given by filename from the given FileIndex.
func LoadRules(filename string, fi FileIndex, eout io.Writer) Rules {
	var rules Rules

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

	rules = ParseRules(data, eout)
	return rules
}

// ParseRules parses the YAML.
func ParseRules(data []byte, eout io.Writer) Rules {
	var rules Rules
	err := yaml.Unmarshal(data, &rules)
	if err != nil {
		fmt.Fprintln(eout, "Error: Could not parse rules file", err)
		return rules
	}

	return rules
}

// ForEachMatchingFile executes a function on each file if the function match returns true.
func ForEachMatchingFile(files []string, rules Rules, out io.Writer, eout io.Writer, fi FileIndex) {
	for _, f := range files {
		matches, rule := FileMatchesRules(f, rules, eout, fi)

		if !matches {
			fmt.Fprintln(out, "No rule for", f)
			continue
		}

		err := ExecuteRule(f, rule, out, eout, fi)
		if err != nil {
			fmt.Fprintf(eout, "Aborting...")
			return
		}

	}
}

// FileMatchesRules if filename matches a rule given in rules returns true and the rule.
func FileMatchesRules(filename string, rules Rules, eout io.Writer, fi FileIndex) (bool, Rule) {
	for _, rule := range rules {
		if FileMatchesRule(filename, rule, eout) {
			return true, rule
		}
	}
	return false, nil
}

// FileMatchesRule return true if filename matches the rule given.
func FileMatchesRule(filename string, rule Rule, eout io.Writer) bool {
	if len(rule["file"]) == 0 {
		fmt.Fprintln(eout, "Error: Empty regexp.")
		return false
	}

	re, err := regexp.Compile(rule["file"])
	if err != nil {
		fmt.Fprintln(eout, "Error: Could not compile regexp", err)
		return false
	}

	return re.MatchString(filename)
}

// ExecuteRule executes the given rules on file.
func ExecuteRule(filename string, rule Rule, out io.Writer, eout io.Writer, fi FileIndex) error {
	dest := rule["move"]

	if !fi.Exists(dest) || !fi.IsDir(dest) {
		fmt.Fprintf(eout, "Error: Invalid directory %s\n", dest)
		return nil
	}
	dest = path.Join(dest, path.Base(filename))

	if fi.Exists(dest) {
		fmt.Fprintln(eout, "Error: File already exists", dest)
		return nil
	}

	fmt.Fprintf(out, "mv -v \"%s\" \"%s\"\n", filename, dest)
	err := fi.Move(filename, dest)

	if err != nil {
		fmt.Fprint(eout, "Error:", err)
		return err
	}

	//err = fi.Remove(filename)
	//if err != nil {
	//fmt.Fprint(eout, "Error:", err)
	//}
	return nil
}
