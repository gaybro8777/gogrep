// Copyright (c) 2017, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package main

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"
)

func loadPaths(wd string, fset *token.FileSet, paths []string) ([]ast.Node, error) {
	var nodes []ast.Node
	ctx := build.Default
	addFile := func(path string) error {
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		nodes = append(nodes, f)
		return nil
	}
	for _, path := range paths {
		if strings.HasSuffix(path, ".go") {
			if err := addFile(path); err != nil {
				return nil, err
			}
			continue
		}
		pkg, err := ctx.Import(path, wd, 0)
		if err != nil {
			return nil, err
		}
		for _, names := range [...][]string{
			pkg.GoFiles, pkg.CgoFiles, pkg.IgnoredGoFiles,
			pkg.TestGoFiles, pkg.XTestGoFiles,
		} {
			for _, name := range names {
				path := filepath.Join(pkg.Dir, name)
				if err := addFile(path); err != nil {
					return nil, err
				}
			}
		}
	}
	return nodes, nil
}

func loadTyped(wd string, fset *token.FileSet, paths []string) (*loader.Program, error) {
	conf := loader.Config{Fset: fset, Cwd: wd}
	if _, err := conf.FromArgs(paths, true); err != nil {
		return nil, err
	}
	return conf.Load()
}