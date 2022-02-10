package errhell

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

func Rewrite(filename string) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	checkErr(err)

	astutil.Apply(f,
		func(c *astutil.Cursor) bool {
			n := c.Node()
			switch v := n.(type) {
			case *ast.FuncDecl:
				handleFuncDecl(v)
			case *ast.FuncLit:
				handleFuncLit(v)
			default:
				return true
			}
			return false
		},
		nil)

	printer.Fprint(os.Stdout, fset, f)
}

func handleFuncDecl(f *ast.FuncDecl) {
	if f.Body == nil {
		return
	}

	returnStack = append(returnStack, f.Type.Results)
	defer func() {
		returnStack = returnStack[:len(returnStack)-1]
	}()

	handleFuncBody(f.Body)
}

func handleFuncLit(f *ast.FuncLit) {
	if f.Body == nil {
		return
	}

	returnStack = append(returnStack, f.Type.Results)
	defer func() {
		returnStack = returnStack[:len(returnStack)-1]
	}()

	handleFuncBody(f.Body)
}

func handleFuncBody(b *ast.BlockStmt) {
	astutil.Apply(b, finder, nil)
}

func checkExprStmt(e *ast.ExprStmt) bool {
	if s, ok := e.X.(*ast.SelectorExpr); ok {
		// TODO support try2, try3 to specify the error position
		if s.Sel.Name == marker {
			astutil.Apply(s.X, finder, nil)
			return true
		}
	}
	return false
}

func finder(c *astutil.Cursor) bool {
	n := c.Node()
	switch v := n.(type) {
	case *ast.FuncLit:
		handleFuncLit(v)
	case *ast.ExprStmt:
		if checkExprStmt(v) {
			rhs := v.X.(*ast.SelectorExpr).X

			assign := &ast.AssignStmt{}
			assign.Tok = token.DEFINE
			assign.Lhs = []ast.Expr{&ast.Ident{Name: errName}}
			assign.Rhs = []ast.Expr{rhs}
			c.Replace(assign)

			ifErr := &ast.IfStmt{}
			ifErr.Cond = &ast.BinaryExpr{
				X:  &ast.Ident{Name: errName},
				Op: token.NEQ,
				Y:  &ast.Ident{Name: "nil"},
			}
			list := make([]ast.Stmt, 0, 2)
			if varStmt := genVar(); varStmt != nil {
				list = append(list, varStmt)
			}
			list = append(list, genReturn())

			ifErr.Body = &ast.BlockStmt{List: list}
			c.InsertAfter(ifErr)
			return true
		}
	case *ast.AssignStmt:
		// TODO
	default:
		return true
	}
	return false
}

func genVar() *ast.DeclStmt {
	currReturn := returnStack[len(returnStack)-1]
	if currReturn == nil {
		return nil
	}

	genDecl := &ast.GenDecl{Tok: token.VAR}
	for i, field := range currReturn.List {
		t := field.Type.(*ast.Ident)
		if t.Name == "error" {
			continue
		}
		valueSpec := &ast.ValueSpec{}
		valueSpec.Names = []*ast.Ident{{Name: fmt.Sprintf("v%d", i)}}
		valueSpec.Type = &ast.Ident{Name: t.Name}
		genDecl.Specs = append(genDecl.Specs, valueSpec)
	}
	if len(genDecl.Specs) == 0 {
		return nil
	}

	ret := &ast.DeclStmt{}
	ret.Decl = genDecl
	return ret
}

func genReturn() *ast.ReturnStmt {
	ret := &ast.ReturnStmt{}
	currReturn := returnStack[len(returnStack)-1]
	if currReturn == nil {
		return ret
	}
	for i, field := range currReturn.List {
		t := field.Type.(*ast.Ident)
		if t.Name == "error" {
			ret.Results = append(ret.Results, &ast.Ident{Name: errName})
			continue
		}
		ret.Results = append(ret.Results, &ast.Ident{Name: fmt.Sprintf("v%d", i)})
	}
	return ret
}
