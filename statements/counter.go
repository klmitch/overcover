// Copyright (c) 2020 Kevin L. Mitchell
//
// Licensed under the Apache License, Version 2.0 (the "License"); you
// may not use this file except in compliance with the License.  You
// may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.  See the License for the specific language governing
// permissions and limitations under the License.

package statements

import (
	"go/ast"
	"path/filepath"

	"golang.org/x/tools/go/packages"

	"github.com/klmitch/overcover/common"
)

// Patch points for top-level functions called by functions in this
// file.
var (
	walk                                                                = ast.Walk
	load func(*packages.Config, ...string) ([]*packages.Package, error) = packages.Load
)

// funcVisitor is a type implementing the ast.Visitor interface.  This
// implementation prospects for function declarations and literals,
// then constructs and returns a stmtVisitor to actually count the
// statements.
type funcVisitor struct {
	fd *common.FileData // The file data
}

// Visit implements the ast.Visitor interface for funcVisitor.
func (v *funcVisitor) Visit(n ast.Node) ast.Visitor {
	switch n.(type) {
	case *ast.FuncDecl:
		return &stmtVisitor{fd: v.fd}
	case *ast.FuncLit:
		return &stmtVisitor{fd: v.fd}
	}

	return v
}

// stmtVisitor is a type implementing the ast.Visitor interface.  This
// implementation prospects for statements and counts them.  Certain
// statements containing other statements are handled specially to
// ensure a proper count, avoiding double-counts.
type stmtVisitor struct {
	fd *common.FileData // The file data
}

// Visit implements the ast.Visitor interface for stmtVisitor.
func (v *stmtVisitor) Visit(n ast.Node) ast.Visitor {
	switch t := n.(type) {
	case *ast.CaseClause: // Handle 'case' in a switch
		for _, stmt := range t.Body {
			walk(v, stmt)
		}
		return nil // we handled recursion ourselves

	case *ast.CommClause: // Handle 'case' in a select
		for _, stmt := range t.Body {
			walk(v, stmt)
		}
		return nil // we handled recursion ourselves

	case *ast.ForStmt: // Handle for statement
		// Count the statement
		v.fd.Count++
		walk(v, t.Body)
		return nil // we handled recursion ourselves

	case *ast.IfStmt: // Handle if statement
		// Count the statement
		v.fd.Count++
		walk(v, t.Body)
		if t.Else != nil {
			walk(v, t.Else)
		}
		return nil // we handled recursion ourselves

	case *ast.SwitchStmt: // Handle switch statement
		// Count the statement
		v.fd.Count++
		walk(v, t.Body)
		return nil // we handled recursion ourselves

	case *ast.TypeSwitchStmt: // Handle type switch statement
		// Count the statement
		v.fd.Count++
		walk(v, t.Body)
		return nil // we handled recursion ourselves

	// Omits BlockStmt and LabeledStmt, as they're not independent
	case *ast.AssignStmt, *ast.BadStmt, *ast.BranchStmt, *ast.DeclStmt, *ast.DeferStmt, *ast.EmptyStmt, *ast.ExprStmt, *ast.GoStmt, *ast.IncDecStmt, *ast.RangeStmt, *ast.ReturnStmt, *ast.SelectStmt, *ast.SendStmt:
		// Count the statement
		v.fd.Count++
	}

	return v
}

// Load loads file data for all files in the specified list of
// packages.  The flags is a list of build flags (nil is acceptable)
// for selecting the files to load, and patterns is a list of package
// patterns to load.  A list of FileData instances is returned.
func Load(flags, patterns []string) (common.DataSet, error) {
	// Begin by constructing the configuration for packages.Load
	cfg := &packages.Config{
		Mode:       packages.NeedTypes | packages.NeedSyntax,
		BuildFlags: flags,
	}

	// Now load the packages
	pkgs, err := load(cfg, patterns...)
	if err != nil {
		return nil, err
	}

	// Next, set up the statement counts
	var data []common.FileData
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			// Assemble the file data
			fd := common.FileData{
				Package: pkg.ID,
				Name:    filepath.Base(pkg.Fset.Position(file.Package).Filename),
			}

			// Count the statements in the file
			walk(&funcVisitor{fd: &fd}, file)

			// Append it to our results
			data = append(data, fd)
		}
	}

	return data, nil
}
