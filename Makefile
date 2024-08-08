GOEXPERIMENT := rangefunc
export GOEXPERIMENT

COVERAGE_TARGET := coverage.txt

.PHONY: build
build:
	go build

.PHONY: run
run: build
	$(TARGET)

.PHONY: test
test:
	go test -coverprofile=$(COVERAGE_TARGET) ./...

.PHONY: format
fmt:
	go mod tidy
	go fmt ./...
	mdformat .