package errhell

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
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

	format.Node(os.Stdout, fset, f)
}

func handleFuncDecl(f *ast.FuncDecl) {
	if f.Body == nil {
		return
	}

	returnStack = append(returnStack, f.Type.Results)
	errVarNameStack = append(errVarNameStack, 0)
	defer func() {
		returnStack = returnStack[:len(returnStack)-1]
		errVarNameStack = errVarNameStack[:len(errVarNameStack)-1]
	}()

	handleFuncBody(f.Body)
}

func handleFuncLit(f *ast.FuncLit) {
	if f.Body == nil {
		return
	}

	returnStack = append(returnStack, f.Type.Results)
	errVarNameStack = append(errVarNameStack, 0)
	defer func() {
		returnStack = returnStack[:len(returnStack)-1]
		errVarNameStack = errVarNameStack[:len(errVarNameStack)-1]
	}()

	handleFuncBody(f.Body)
}

func handleFuncBody(b *ast.BlockStmt) {
	astutil.Apply(b, finder, nil)
}

// checkExprStmt will return integer N representing the Nth return value should
// be used to build if-err.
func checkExprStmt(e *ast.ExprStmt) int {
	if s, ok := e.X.(*ast.SelectorExpr); ok {
		return checkSelectorExpr(s)
	}
	return 0
}

// checkSelectorExpr will return integer N representing the Nth return value
// should be used to build if-err.
func checkSelectorExpr(s *ast.SelectorExpr) int {
	if i := matchMarker(s.Sel.Name); i > 0 {
		astutil.Apply(s.X, finder, nil)
		return i
	}
	return 0
}

func finder(c *astutil.Cursor) bool {
	n := c.Node()
	switch v := n.(type) {
	case *ast.FuncLit:
		handleFuncLit(v)
	case *ast.ExprStmt:
		// TODO: support nesting this ExprStmt in a if
		i := checkExprStmt(v)
		if i == 0 {
			break
		}
		lhs := []ast.Expr{}
		for ; i > 1; i-- {
			lhs = append(lhs, &ast.Ident{Name: "_"})
		}
		errVarName := genErrVarName()
		lhs = append(lhs, &ast.Ident{Name: errVarName})
		rhs := []ast.Expr{v.X.(*ast.SelectorExpr).X}

		assign := &ast.AssignStmt{}
		assign.Tok = token.DEFINE
		assign.Lhs = lhs
		assign.Rhs = rhs
		c.Replace(assign)

		c.InsertAfter(genIfErr(errVarName))
		return true
	case *ast.AssignStmt:
		if len(v.Rhs) != 1 {
			break
		}
		s, ok := v.Rhs[0].(*ast.SelectorExpr)
		if !ok {
			break
		}

		i := checkSelectorExpr(s)
		if i == 0 {
			break
		}

		v.Rhs[0] = s.X
		c.Replace(v)
		c.InsertAfter(genIfErr(v.Lhs[i-1].(*ast.Ident).Name))
		return true
	default:
		return true
	}
	return false
}

func genIfErr(errVarName string) *ast.IfStmt {
	ifErr := &ast.IfStmt{}
	ifErr.Cond = &ast.BinaryExpr{
		X:  &ast.Ident{Name: errVarName},
		Op: token.NEQ,
		Y:  &ast.Ident{Name: "nil"},
	}
	list := make([]ast.Stmt, 0, 2)
	// TODO: for primitive type, use literal instead declare a variable
	if varStmt := genVar(); varStmt != nil {
		list = append(list, varStmt)
	}
	list = append(list, genReturn(errVarName))

	ifErr.Body = &ast.BlockStmt{List: list}
	return ifErr
}

func genVar() *ast.DeclStmt {
	currReturn := returnStack[len(returnStack)-1]
	if currReturn == nil {
		return nil
	}

	genDecl := &ast.GenDecl{Tok: token.VAR}
	for i, tp := range extractTypes(currReturn.List) {
		if t, ok := tp.(*ast.Ident); ok {
			if t.Name == "error" {
				continue
			}
		}
		if zeroValueLiteralForType(tp) != nil {
			continue
		}

		valueSpec := &ast.ValueSpec{}
		valueSpec.Names = []*ast.Ident{{Name: fmt.Sprintf("_v%d", i)}}
		valueSpec.Type = tp
		genDecl.Specs = append(genDecl.Specs, valueSpec)
	}
	if len(genDecl.Specs) == 0 {
		return nil
	}

	ret := &ast.DeclStmt{}
	ret.Decl = genDecl
	return ret
}

func genReturn(errVarName string) *ast.ReturnStmt {
	ret := &ast.ReturnStmt{}
	currReturn := returnStack[len(returnStack)-1]
	if currReturn == nil {
		return ret
	}
	for i, tp := range extractTypes(currReturn.List) {
		if t, ok := tp.(*ast.Ident); ok {
			if t.Name == "error" {
				ret.Results = append(ret.Results, &ast.Ident{Name: errVarName})
				continue
			}
		}
		if expr := zeroValueLiteralForType(tp); expr != nil {
			ret.Results = append(ret.Results, expr)
			continue
		}

		ret.Results = append(ret.Results, &ast.Ident{Name: fmt.Sprintf("_v%d", i)})
	}
	return ret
}
