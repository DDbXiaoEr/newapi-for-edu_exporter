.PHONY: build run clean test

BINARY_NAME=newapi-exporter
BUILD_DIR=bin

build:
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) .

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test ./... -v

clean:
	rm -rf $(BUILD_DIR) output

lint:
	go vet ./...
