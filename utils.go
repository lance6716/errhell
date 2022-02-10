package errhell

import (
	"fmt"
	"go/ast"
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
