package main

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestMain_Simple2(t *testing.T) {
	src := `
	package assert

import (
        "unicode"
        "unicode/utf8"
        "github.com/davecgh/go-spew/spew"
        "github.com/pmezard/go-difflib/difflib"

	fun main() {}
	`
	fs := token.NewFileSet()
	f, _ := parser.ParseFile(fs, "", src, parser.ImportsOnly)
	updated := update("", f, "unicode", "github.org/encode")
	if !updated {
		t.FailNow()
	}
	found := false
	for _, i := range f.Imports {
		found = strings.Index(i.Path.Value, "unicode") != -1 || found
	}

	if found {
		t.FailNow()
	}
}

func TestMain_Simple(t *testing.T) {
	src := `
	package assert

import (
        "unicode"
        "unicode/utf8"
        "github.com/davecgh/go-spew/spew"
        "github.com/pmezard/go-difflib/difflib"

	fun main() {}
	`
	fs := token.NewFileSet()
	f, _ := parser.ParseFile(fs, "", src, parser.ImportsOnly)
	updated := update("", f, "github.com", "github.org")
	if !updated {
		t.FailNow()
	}
	found := false
	for _, i := range f.Imports {
		found = strings.Index(i.Path.Value, "github.com") != -1 || found
	}

	if found {
		t.FailNow()
	}
}
