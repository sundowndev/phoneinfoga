package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2/internal/takeon/github.com/markbates/errx"
)

// ParsedFile ...
type ParsedFile struct {
	File    packd.SimpleFile
	FileSet *token.FileSet
	Ast     *ast.File
	Lines   []string
}

// ParseFileMode ...
func ParseFileMode(gf packd.SimpleFile, mode parser.Mode) (ParsedFile, error) {
	pf := ParsedFile{
		FileSet: token.NewFileSet(),
		File:    gf,
	}

	src := gf.String()
	f, err := parser.ParseFile(pf.FileSet, gf.Name(), src, mode)
	if err != nil && errx.Unwrap(err) != io.EOF {
		return pf, err
	}
	pf.Ast = f

	pf.Lines = strings.Split(src, "\n")
	return pf, nil
}

// ParseFile ...
func ParseFile(gf packd.SimpleFile) (ParsedFile, error) {
	return ParseFileMode(gf, 0)
}
