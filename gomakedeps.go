package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

type Target struct {
	Binary  string
	Imports []*Import
}

type Import struct {
	Path    string
	PkgPath string
	Goroot  bool
}

func lookupDependencies(filePath string) ([]*Import, error) {
	// resolve filepath
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		panic(err)
	}
	fileDir := filepath.Dir(absFilePath)
	pkg, err := build.ImportDir(fileDir, 0)
	if err != nil {
		panic(err)
	}

	imports := make([]*Import, 0, len(pkg.ImportPos))

	for importName, importPos := range pkg.ImportPos {
		// load package for import
		for _, pos := range importPos {
			if pos.Filename == absFilePath {
				// load package for import to resolve the absolute import path
				importPkg, err := build.Import(importName, fileDir, 0)
				if err != nil {
					panic(err)
				}
				importPath := fmt.Sprintf("%s/src/%s", importPkg.Root, importPkg.ImportPath)
				imports = append(imports, &Import{Path: importName, PkgPath: importPath, Goroot: importPkg.Goroot})
			}
		}
	}
	return imports, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <go main sourcefile>\n", os.Args[0])
		os.Exit(1)
	}
	filePath := os.Args[1]
	if imports, err := lookupDependencies(filePath); err != nil {
		fmt.Fprintln(os.Stderr, "Depency lookup failed for:", filePath, err.Error())
	} else {
		for _, i := range imports {
			fmt.Print(i.PkgPath + "/*.go ")
		}
	}
}
