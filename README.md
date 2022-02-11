# errhell

A preprocessor for golang to generate notorious `if err != nil {...}`.

## Example

Input:

```golang
package main

func main() {
	returnErr().try
}

func returnErr() error {
	return nil
}

```

```bash
bin/errhell ./test.go
```

Output:

```golang
package main

func main() {
	_err0 := returnErr()
	if _err0 != nil {
		return
	}
}

func returnErr() error {
	return nil
}

```

You can find more example in `test/input/input.go`.
