package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"text/template"
	"strings"
)

type Target struct {
	Binary          string
	Imports					[]*Import
}

type Import struct {
	Path    string
	PkgPath string
	Goroot  bool
}

const makefileTemplate =
`{{ .Binary }}: {{ range $i := .Imports }}{{ $i.PkgPath }}/*.go{{ end }}
	go build -ldflags ${LDFLAGS} $@.go
`

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
	// TODO add flag to show only import from GOPATH (not GOROOT)
	var templatePath string
	flag.StringVar(&templatePath, "template", "", "Makefile template path.")
	flag.Parse()

	// parse template
	var t *template.Template
	if templatePath == "" {
		t = template.Must(template.New("Makefile").Parse(makefileTemplate))
	} else {
		t = template.Must(template.ParseFiles(templatePath))
	}


	filePath := flag.Args()[0]
	if imports, err := lookupDependencies(filePath); err != nil {
		fmt.Fprintln(os.Stderr, "Depency lookup failed for:", filePath, err.Error())
	} else {
		target := &Target{
			Binary: strings.TrimSuffix(filepath.Base(filePath), ".go"),
			Imports: imports,
		}
		// render template
		t.Execute(os.Stdout, target)
	}
}

