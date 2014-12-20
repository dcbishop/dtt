package main

import (
	"bytes"
	"log"
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

func TestMainMissingFile(t *testing.T) {
	missingFilename := "thisfileshouldntexist_asd3f4f2tfdsfa"
	out := bytes.NewBufferString("")

	lf := NewLocalFiles("")
	Main([]string{"dealwithit", missingFilename}, out, out, &lf)

	if out.String() != "Error: File not found: "+missingFilename+"\n" {
		log.Println(out.String())
		t.Error("Did not print usage.")
	}
}
