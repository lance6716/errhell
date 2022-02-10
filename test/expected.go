package main

import "fmt"

func main() {
	_err0 := foo()
	if _err0 != nil {
		return
	}
}

func foo() error {
	i, myErr := a()
	if myErr != nil {
		return myErr
	}
	i++
	return nil
}

func bar() (i int, err error) {
	err = b(0)
	if err != nil {
		var v0 int
		return v0, err
	}
	return 0, nil
}

func b(int) error {
	return nil
}

var a = func() (float32, error) {
	_, _err0 := fmt.Printf("%p", b)
	if _err0 != nil {
		var v0 float32
		return v0, _err0
	}
	_, _, _err1 := func() (string, string, error) {
		var myErr error
		myErr = foo()
		if myErr != nil {
			var (
				v0 string
				v1 string
			)
			return v0, v1, myErr
		}
		return "", "", nil
	}()
	if _err1 != nil {
		var v0 float32
		return v0, _err1
	}

	return 0, nil
}
