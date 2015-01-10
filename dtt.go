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
	"github.com/spf13/afero"
	"github.com/termie/go-shutil"
	"gopkg.in/yaml.v1"
)

func main() {
	Main(os.Args, os.Stdout, os.Stderr, &afero.OsFs{})
}

// The programs usage help.
var Usage = `dtt.

Usage:
	dtt <file>
`

// Exists returns true if a file given by path exists.
func Exists(fs afero.Fs, filename string) bool {
	if _, err := fs.Stat(filename); err != nil {
		return false
	}
	return true
}

// IsDir returns true if filename points to a directory.
func IsDir(fs afero.Fs, filename string) bool {
	fileInfo, err := fs.Stat(filename)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// Move moves a directory or file to the destination directory.
func Move(fs afero.OsFs, oldpath, newpath string) error {
	cmd := exec.Command("mv", "-v", oldpath, newpath)
	out, err := cmd.Output()
	fmt.Println(string(out))
	return err
}

// CopyWithTemp copies a directory to the destination.
func CopyWithTemp(oldpath, newpath string) error {
	// [TODO]: Actually do temp stuff - 2014-12-21 01:49pm
	src := oldpath
	dst := newpath
	//dstTemp := dst + ".diw-COPYING"
	//err := shutil.CopyTree(src, dstTemp, nil)
	err := shutil.CopyTree(src, dst, nil)
	if err != nil {
		return err
	}
	return err

	//return os.Rename(dstTemp, dst)
}

// Main function.
func Main(args []string, out io.Writer, eout io.Writer, fi afero.Fs) {
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

// AnyFilesDontExist will look for each of the given files in the given afero.Fs,
// output an error to eout for any file that doesn't exist and return true if where any missing.
func AnyFilesDontExist(eout io.Writer, fi afero.Fs, files []string) bool {
	err := false
	for _, f := range files {
		if !Exists(fi, f) {
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

// LoadRules loads the rules file given by filename from the given afero.Fs.
func LoadRules(filename string, fi afero.Fs, eout io.Writer) Rules {
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
func ForEachMatchingFile(files []string, rules Rules, out io.Writer, eout io.Writer, fi afero.Fs) {
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
func FileMatchesRules(filename string, rules Rules, eout io.Writer, fi afero.Fs) (bool, Rule) {
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
func ExecuteRule(filename string, rule Rule, out io.Writer, eout io.Writer, fi afero.Fs) error {
	dest := rule["move"]

	if !Exists(fi, dest) || !IsDir(fi, dest) {
		fmt.Fprintf(eout, "Error: Invalid directory %s\n", dest)
		return nil
	}
	dest = path.Join(dest, path.Base(filename))

	if Exists(fi, dest) {
		fmt.Fprintln(eout, "Error: File already exists", dest)
		return nil
	}

	fmt.Fprintf(out, "mv -v \"%s\" \"%s\"\n", filename, dest)
	osfs := fi.(*afero.OsFs)
	err := Move(*osfs, filename, dest)

	if err != nil {
		fmt.Fprint(eout, "Error:", err)
		return err
	}

	return nil
}
