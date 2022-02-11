package errhell

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

var marker = "try"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var returnStack []*ast.FieldList
var errVarNameStack []int

func genErrVarName() string {
	curr := errVarNameStack[len(errVarNameStack)-1]
	ret := fmt.Sprintf("_err%d", curr)
	curr++
	errVarNameStack[len(errVarNameStack)-1] = curr
	return ret
}

// matchMarker will extract the number N which in the pattern $(marker)N.
// empty N returns 1.
func matchMarker(name string) int {
	if name == marker {
		return 1
	}
	if !strings.HasPrefix(name, marker) {
		return 0
	}
	i, err := strconv.Atoi(name[len(marker):])
	if err != nil {
		return 0
	}
	return i
}

// extractTypes flattens multiple names belongs to one type.
func extractTypes(list []*ast.Field) []ast.Expr {
	ret := []ast.Expr{}
	for _, f := range list {
		if len(f.Names) == 0 {
			ret = append(ret, f.Type)
			continue
		}
		for range f.Names {
			ret = append(ret, f.Type)
		}
	}
	return ret
}

var zeroValueForBasicType = map[string]ast.Expr{
	"bool":       &ast.Ident{Name: "false"},
	"string":     &ast.BasicLit{Kind: token.STRING, Value: `""`},
	"int":        &ast.BasicLit{Kind: token.INT, Value: "0"},
	"int8":       &ast.BasicLit{Kind: token.INT, Value: "0"},
	"int16":      &ast.BasicLit{Kind: token.INT, Value: "0"},
	"int32":      &ast.BasicLit{Kind: token.INT, Value: "0"},
	"int64":      &ast.BasicLit{Kind: token.INT, Value: "0"},
	"uint":       &ast.BasicLit{Kind: token.INT, Value: "0"},
	"uint8":      &ast.BasicLit{Kind: token.INT, Value: "0"},
	"uint16":     &ast.BasicLit{Kind: token.INT, Value: "0"},
	"uint32":     &ast.BasicLit{Kind: token.INT, Value: "0"},
	"uint64":     &ast.BasicLit{Kind: token.INT, Value: "0"},
	"uintptr":    &ast.BasicLit{Kind: token.INT, Value: "0"},
	"byte":       &ast.BasicLit{Kind: token.INT, Value: "0"},
	"rune":       &ast.BasicLit{Kind: token.INT, Value: "0"},
	"float32":    &ast.BasicLit{Kind: token.INT, Value: "0"},
	"float64":    &ast.BasicLit{Kind: token.INT, Value: "0"},
	"complex64":  &ast.BasicLit{Kind: token.INT, Value: "0"},
	"complex128": &ast.BasicLit{Kind: token.INT, Value: "0"},
}

// zeroValueLiteralForType returns an expression of zero value of the given type
// which could be written as valid source code. Return nil if not support this
// type.
func zeroValueLiteralForType(tp ast.Expr) ast.Expr {
	switch v := tp.(type) {
	case *ast.Ident:
		return zeroValueForBasicType[v.Name]
	case *ast.StarExpr, *ast.FuncType, *ast.ChanType, *ast.ArrayType,
		*ast.MapType, *ast.InterfaceType:
		return &ast.Ident{Name: "nil"}
	}
	return nil
}
