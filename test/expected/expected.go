package test

import (
	"fmt"
	"io"
)

func main() {
	_err0 := returnError()
	if _err0 != nil {
		return
	}
}

func returnError() error {
	i, myErr := funcVarReturnFloat32Error()
	if myErr != nil {
		return myErr
	}
	i++
	_ = myErr
	return nil
}

func returnIntError() (i int, err error) {
	err = funcWithParam(0)
	if err != nil {
		var v0 int
		return v0, err
	}

	_, err = returnReaderError()
	if err != nil {
		var v0 int
		return v0, err
	}

	return 0, nil
}

func funcWithParam(int) error {
	return nil
}

var funcVarReturnFloat32Error = func() (float32, error) {
	_, _err0 := fmt.Printf("%p", funcWithParam)
	if _err0 != nil {
		var v0 float32
		return v0, _err0
	}
	_, _, _err1 := func() (string, string, error) {
		var myErr error
		myErr = returnError()
		if myErr != nil {
			var (
				v0 string
				v1 string
			)
			return v0, v1, myErr
		}

		_ = myErr
		return "", "", nil
	}()
	if _err1 != nil {
		var v0 float32
		return v0, _err1
	}

	return 0, nil
}

func returnReaderError() (io.Reader, error) {
	_err0 := returnError()
	if _err0 != nil {
		var v0 io.Reader
		return v0, _err0
	}

	return nil, nil
}

func returnInt() int {
	_err0 := returnError()
	if _err0 != nil {
		var v0 int
		return v0
	}

	return 0
}

func namedReturn() (a, b int, err error) {
	_err0 := returnError()
	if _err0 != nil {
		var (
			v0 int
			v1 int
		)
		return v0, v1, _err0
	}

	return 0, 0, err
}

func manyReturn() (complex64,
	*int,
	<-chan struct{},
	[]float32,
	interface{},
	func(error) error,
	map[string]interface{},
	error) {
	_err0 := returnError()
	if _err0 != nil {
		var (
			v0 complex64
			v1 *int
			v2 <-chan struct{}
			v3 []float32
			v4 interface{}
			v5 func(error) error
			v6 map[string]interface{}
		)
		return v0, v1, v2, v3, v4, v5, v6, _err0
	}

	return 0, nil, nil, nil, nil, nil, nil, nil
}
