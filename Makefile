MODULE_DIR := usm
TARGET := build/usm

.PHONY: build
build:
	go build -o $(TARGET) $(MODULE_DIR)

.PHONY: run
run: build
	$(TARGET)

.PHONY: test
test:
	go test $(MODULE_DIR)/...