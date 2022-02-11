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
		return 0, err
	}
	_, err = returnReaderError()
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func funcWithParam(int) error {
	return nil
}

var funcVarReturnFloat32Error = func() (float32, error) {
	_, _err0 := fmt.Printf("%p", funcWithParam)
	if _err0 != nil {
		return 0, _err0
	}
	_, _, _err1 := func() (string, string, error) {
		var myErr error
		myErr = returnError()
		if myErr != nil {
			return "", "", myErr
		}
		_ = myErr
		return "", "", nil
	}()
	if _err1 != nil {
		return 0, _err1
	}

	return 0, nil
}

func returnReaderError() (io.Reader, error) {
	_err0 := returnError()
	if _err0 != nil {
		var _v0 io.Reader
		return _v0, _err0
	}

	return nil, nil
}

func returnInt() int {
	_err0 := returnError()
	if _err0 != nil {
		return 0
	}
	return 0
}

func namedReturn() (a, b int, err error) {
	_err0 := returnError()
	if _err0 != nil {
		return 0, 0, _err0
	}
	return 0, 0, err
}

func manyReturn() (
	complex64,
	*int,
	<-chan struct{},
	[]float32,
	interface{},
	func(error) error,
	map[string]interface{},
	error,
) {
	_err0 := returnError()
	if _err0 != nil {
		return 0, nil, nil, nil, nil, nil, nil, _err0
	}
	return 0, nil, nil, nil, nil, nil, nil, nil
}
