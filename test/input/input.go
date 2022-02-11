package test

import (
	"fmt"
	"io"
)

func main() {
	returnError().try
}

func returnError() error {
	i, myErr := funcVarReturnFloat32Error().try2
	i++
	_ = myErr
	return nil
}

func returnIntError() (i int, err error) {
	err = funcWithParam(0).try
	_, err = returnReaderError().try2
	return 0, nil
}

func funcWithParam(int) error {
	return nil
}

var funcVarReturnFloat32Error = func() (float32, error) {
	fmt.Printf("%p", funcWithParam).try2

	func() (string, string, error) {
		var myErr error
		myErr = returnError().try
		_ = myErr
		return "", "", nil
	}().try3

	return 0, nil
}

func returnReaderError() (io.Reader, error) {
	returnError().try
	return nil, nil
}

func returnInt() int {
	returnError().try
	return 0
}

func namedReturn() (a, b int, err error) {
	returnError().try
	return 0, 0, err
}
