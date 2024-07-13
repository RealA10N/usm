MODULE_DIR := usm
TARGET := build/usm
COVERAGE_TARGET := coverage.out

.PHONY: build
build:
	go build -o $(TARGET) $(MODULE_DIR)

.PHONY: run
run: build
	$(TARGET)

.PHONY: test
test:
	go test -coverprofile=$(COVERAGE_TARGET) $(MODULE_DIR)/...