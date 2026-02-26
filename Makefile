.PHONY: build run test test-coverage clean

BINARY=flaggyd
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/flaggyd

run: build
	$(BUILD_DIR)/$(BINARY)

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./internal/engine/...
	go tool cover -func=coverage.out | tail -1
	@echo "---"
	@echo "Full coverage report: go tool cover -html=coverage.out"

clean:
	rm -rf $(BUILD_DIR) coverage.out flaggy.db
