package main

import (
	"bytes"
	"testing"
)

func TestMainNoArgsShouldPrintUsage(t *testing.T) {
	out := bytes.NewBufferString("")

	lf := NewLocalFiles("")
	Main([]string{""}, out, out, &lf)

	if out.String() != Usage+"\n" {
		t.Error("Did not print usage.")
	}
}
