TARGET := build/usm
COVERAGE_TARGET := coverage.txt

GOEXPERIMENT := rangefunc
export GOEXPERIMENT

.PHONY: build
build:
	go build -o $(TARGET)

.PHONY: run
run: build
	$(TARGET)

.PHONY: test
test:
	go test -coverprofile=$(COVERAGE_TARGET) ./...
