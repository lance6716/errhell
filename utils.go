package errhell

import (
	"go/ast"
)

var marker = "try"
var errName = "_err_"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var returnStack []*ast.FieldList
