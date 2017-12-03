package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	packagePath = flag.String("package", "", "the package path")
	from        = flag.String("from", "", "the old import that will be replaced")
	replaceWith = flag.String("to", "", "the new import")
	dryRun      = flag.String("dryRun", "false", "prints out the changes. Does not edit files")
)

func main() {
	var err error
	flag.Parse()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, *packagePath, nil, parser.ImportsOnly)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, p := range pkgs {
		for fname, f := range p.Files {
			if update(fname, f, *from, *replaceWith) {
				fop, err := os.OpenFile(fname, os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				err = printer.Fprint(fop, fset, f)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				ast.SortImports(fset, f)
				fmt.Println("File updated:", fname)
			}
		}
	}
}

func update(filename string, file *ast.File, from, to string) bool {
	writeChange := false
	for _, i := range file.Imports {
		val, _ := strconv.Unquote(i.Path.Value)
		if strings.Index(val, from) != -1 {
			replacement := path.Clean(strings.Replace(val, from, to, -1))
			if *dryRun == "false" {
				writeChange = true
				i.Path.Value = strconv.Quote(replacement)
			} else {
				fmt.Println(filename, ":", val, "->", replacement)
			}
		}
	}
	return writeChange
}
