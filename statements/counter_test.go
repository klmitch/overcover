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
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"

	"github.com/klmitch/patcher"

	"github.com/klmitch/overcover/common"
)

func TestFuncVisitorVisitFuncDecl(t *testing.T) {
	obj := &funcVisitor{
		fd: &common.FileData{},
	}
	node := &ast.FuncDecl{}

	result := obj.Visit(node)

	assert.Equal(t, &stmtVisitor{
		fd: obj.fd,
	}, result)
	assert.Same(t, result.(*stmtVisitor).fd, obj.fd)
}

func TestFuncVisitorVisitFuncLit(t *testing.T) {
	obj := &funcVisitor{
		fd: &common.FileData{},
	}
	node := &ast.FuncLit{}

	result := obj.Visit(node)

	assert.Equal(t, &stmtVisitor{
		fd: obj.fd,
	}, result)
	assert.Same(t, result.(*stmtVisitor).fd, obj.fd)
}

func TestFuncVisitorVisitOther(t *testing.T) {
	obj := &funcVisitor{
		fd: &common.FileData{},
	}
	node := &ast.EmptyStmt{}

	result := obj.Visit(node)

	assert.Same(t, obj, result)
}

func TestStmtVisitorVisitCaseClause(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	stmts := []ast.Stmt{
		&ast.EmptyStmt{},
		&ast.EmptyStmt{},
		&ast.EmptyStmt{},
	}
	node := &ast.CaseClause{
		Body: stmts,
	}
	walkCalled := 0
	defer patcher.SetVar(&walk, func(v ast.Visitor, n ast.Node) {
		assert.Same(t, obj, v)
		assert.Same(t, stmts[walkCalled], n)
		walkCalled++
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Nil(t, result)
	assert.Equal(t, int64(0), obj.fd.Count)
	assert.Equal(t, 3, walkCalled)
}

func TestStmtVisitorVisitCommClause(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	stmts := []ast.Stmt{
		&ast.EmptyStmt{},
		&ast.EmptyStmt{},
		&ast.EmptyStmt{},
	}
	node := &ast.CommClause{
		Body: stmts,
	}
	walkCalled := 0
	defer patcher.SetVar(&walk, func(v ast.Visitor, n ast.Node) {
		assert.Same(t, obj, v)
		assert.Same(t, stmts[walkCalled], n)
		walkCalled++
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Nil(t, result)
	assert.Equal(t, int64(0), obj.fd.Count)
	assert.Equal(t, 3, walkCalled)
}

func TestStmtVisitorVisitForStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	body := &ast.BlockStmt{}
	node := &ast.ForStmt{
		Body: body,
	}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, n ast.Node) {
		assert.Same(t, obj, v)
		assert.Same(t, body, n)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Nil(t, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.True(t, walkCalled)
}

func TestStmtVisitorVisitIfStmtBase(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	body := &ast.BlockStmt{}
	node := &ast.IfStmt{
		Body: body,
	}
	walkCalled := 0
	defer patcher.SetVar(&walk, func(v ast.Visitor, n ast.Node) {
		assert.Same(t, obj, v)
		switch walkCalled {
		case 0:
			assert.Same(t, body, n)
		default:
			t.Error("Unexpected call to walk")
			t.Fail()
		}
		walkCalled++
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Nil(t, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.Equal(t, 1, walkCalled)
}

func TestStmtVisitorVisitIfStmtElse(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	body := &ast.BlockStmt{}
	elseStmt := &ast.EmptyStmt{}
	node := &ast.IfStmt{
		Body: body,
		Else: elseStmt,
	}
	walkCalled := 0
	defer patcher.SetVar(&walk, func(v ast.Visitor, n ast.Node) {
		assert.Same(t, obj, v)
		switch walkCalled {
		case 0:
			assert.Same(t, body, n)
		case 1:
			assert.Same(t, elseStmt, n)
		default:
			t.Error("Unexpected call to walk")
			t.Fail()
		}
		walkCalled++
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Nil(t, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.Equal(t, 2, walkCalled)
}

func TestStmtVisitorVisitSwitchStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	body := &ast.BlockStmt{}
	node := &ast.SwitchStmt{
		Body: body,
	}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, n ast.Node) {
		assert.Same(t, obj, v)
		assert.Same(t, body, n)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Nil(t, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.True(t, walkCalled)
}

func TestStmtVisitorVisitTypeSwitchStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	body := &ast.BlockStmt{}
	node := &ast.TypeSwitchStmt{
		Body: body,
	}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, n ast.Node) {
		assert.Same(t, obj, v)
		assert.Same(t, body, n)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Nil(t, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.True(t, walkCalled)
}

func TestStmtVisitorVisitAssignStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.AssignStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitBadStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.BadStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitBranchStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.BranchStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitDeclStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.DeclStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitDeferStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.DeferStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitEmptyStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.EmptyStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitExprStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.ExprStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitGoStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.GoStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitIncDecStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.IncDecStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitRangeStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.RangeStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitReturnStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.ReturnStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitSelectStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.SelectStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestStmtVisitorVisitSendStmt(t *testing.T) {
	obj := &stmtVisitor{
		fd: &common.FileData{},
	}
	node := &ast.SendStmt{}
	walkCalled := false
	defer patcher.SetVar(&walk, func(v ast.Visitor, _ ast.Node) {
		assert.Same(t, obj, v)
		walkCalled = true
	}).Install().Restore()

	result := obj.Visit(node)

	assert.Same(t, obj, result)
	assert.Equal(t, int64(1), obj.fd.Count)
	assert.False(t, walkCalled)
}

func TestLoadBase(t *testing.T) {
	fset := &token.FileSet{}
	p1f1 := fset.AddFile("some/path/p1f1", -1, 4)
	p1f2 := fset.AddFile("some/path/p1f2", -1, 4)
	p2f1 := fset.AddFile("some/path/p2f1", -1, 4)
	nodes := []*ast.File{
		{
			Package: token.Pos(p1f1.Base() + 1),
		},
		{
			Package: token.Pos(p1f2.Base() + 1),
		},
		{
			Package: token.Pos(p2f1.Base() + 1),
		},
	}
	pkgs := []*packages.Package{
		{
			ID:     "p1",
			Fset:   fset,
			Syntax: []*ast.File{nodes[0], nodes[1]},
		},
		{
			ID:     "p2",
			Fset:   fset,
			Syntax: []*ast.File{nodes[2]},
		},
	}
	ds := common.DataSet{
		common.FileData{
			Package: "p1",
			Name:    "p1f1",
		},
		common.FileData{
			Package: "p1",
			Name:    "p1f2",
		},
		common.FileData{
			Package: "p2",
			Name:    "p2f1",
		},
	}
	loadCalled := false
	walkCalled := 0
	defer patcher.NewPatchMaster(
		patcher.SetVar(&load, func(cfg *packages.Config, patterns ...string) ([]*packages.Package, error) {
			assert.Equal(t, &packages.Config{
				Mode:       packages.NeedTypes | packages.NeedSyntax,
				BuildFlags: []string{},
			}, cfg)
			assert.Equal(t, []string{"./..."}, patterns)
			loadCalled = true
			return pkgs, nil
		}),
		patcher.SetVar(&walk, func(v ast.Visitor, n ast.Node) {
			fv, ok := v.(*funcVisitor)
			require.True(t, ok)
			assert.Equal(t, &ds[walkCalled], fv.fd)
			assert.Same(t, nodes[walkCalled], n)
			walkCalled++
		}),
	).Install().Restore()

	result, err := Load([]string{}, []string{"./..."})

	assert.NoError(t, err)
	assert.Equal(t, ds, result)
	assert.True(t, loadCalled)
	assert.Equal(t, 3, walkCalled)
}

func TestLoadError(t *testing.T) {
	loadCalled := false
	walkCalled := 0
	defer patcher.NewPatchMaster(
		patcher.SetVar(&load, func(cfg *packages.Config, patterns ...string) ([]*packages.Package, error) {
			assert.Equal(t, &packages.Config{
				Mode:       packages.NeedTypes | packages.NeedSyntax,
				BuildFlags: []string{},
			}, cfg)
			assert.Equal(t, []string{"./..."}, patterns)
			loadCalled = true
			return nil, assert.AnError
		}),
		patcher.SetVar(&walk, func(_ ast.Visitor, _ ast.Node) {
			walkCalled++
		}),
	).Install().Restore()

	result, err := Load([]string{}, []string{"./..."})

	assert.Same(t, assert.AnError, err)
	assert.Nil(t, result)
	assert.True(t, loadCalled)
	assert.Equal(t, 0, walkCalled)
}
