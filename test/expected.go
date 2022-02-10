package main

import "fmt"

func main() {
	_err_ := foo()
	if _err_ != nil {
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

func bar() (int, error) {
	_err_ := b(0)
	if _err_ != nil {
		var v0 int
		return v0, _err_
	}
	return 0, nil
}

func b(int) error {
	return nil
}

var a = func() (float32, error) {
	_, _err_ := fmt.Printf("%p", b)
	if _err_ != nil {
		var v0 float32
		return v0, _err_
	}
	return 0, nil
}
