package errhell

import (
	"go/ast"
	"strconv"
	"strings"
)

var marker = "try"
var errName = "_err_"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var returnStack []*ast.FieldList

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
