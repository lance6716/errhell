package main

func main() {
	_err_ := foo()
	if _err_ != nil {
		return
	}
}

func foo() error {
	_err_ := b(0)
	if _err_ != nil {
		return _err_
	}
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

var a = func() (int, error) {
	return 0, nil
}
