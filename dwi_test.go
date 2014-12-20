package main

import (
	"bytes"
	"testing"
)

func TestMainNoArgsShouldPrintUsage(t *testing.T) {
	out := bytes.NewBufferString("")
	err := bytes.NewBufferString("")

	Main([]string{""}, out, err, NewLocalFiles("/"))

	if out.String() != Usage+"\n" {
		t.Error("Did not print usage.")
	}
}
