GOEXPERIMENT := rangefunc
export GOEXPERIMENT

.PHONY: build
build:
	go build

.PHONY: run
run: build
	$(TARGET)

.PHONY: test
test:
	go test -v ./...

.PHONY: ci
ci:
	go install github.com/dave/courtney@v0.4.1

.PHONY: coverage
coverage:
	courtney -v

.PHONY: format
fmt:
	go mod tidy
	go fmt ./...
	mdformat .