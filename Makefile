.PHONY: build test
build:
	go build -o bin/errhell cmd/main.go
test: build
	bin/errhell test/input/input.go > /tmp/simple.go
	diff /tmp/simple.go test/expected/expected.go
print-test: build
	bin/errhell test/input/input.go