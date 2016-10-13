package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	//"flag"
	//"text/template"
)

type Target struct {
	Binary string
	DependencyFiles []string
}

type Import struct {
	Path string
	PkgPath string
	Goroot bool
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
	//var templatePath string
	//flag.String("template", &templatePath, "makefile.tmpl", "Makefile template path.")
	//flag.Parse()

	//mainFiles := flag.Args()
	mainFiles := os.Args[1:]
	// TODO add flag to show only import from GOPATH (not GOROOT)

	for _, filePath := range mainFiles {
		if imports, err := lookupDependencies(filePath); err != nil {
			fmt.Fprintln(os.Stderr, "Depency lookup failed for:", filePath, err.Error())
		} else {
			fmt.Println(filePath)
			for _, i := range imports {
				fmt.Printf("%#v\n", i)
			}
		}
	}
}
