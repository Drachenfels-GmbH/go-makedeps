package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

func main() {
	filePath := os.Args[1]
	// TODO add flag to show only import from GOPATH (not GOROOT)

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

	//imports := make(map[string]*build.Package)

	for importName, importPos := range pkg.ImportPos {
		// load package for import
		for _, pos := range importPos {
			if pos.Filename == absFilePath {
				// load package for import to resolve the absolute import path
				importPkg, err := build.Import(importName, fileDir, 0)
				if err != nil {
					panic(err)
				}
				if importPkg.Goroot {
					fmt.Printf("%s/%s/*.go ", build.Default.GOROOT, importPkg.ImportPath)
				} else {
					fmt.Printf("%s/%s/*.go ", build.Default.GOPATH, importPkg.ImportPath)
				}
			}
		}
	}
}
