.PHONY: test
test:
	go run ./cmd/main.go test/simple.go.input > /tmp/simple.go
	diff /tmp/simple.go test/expected.go