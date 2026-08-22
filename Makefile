TARGET_DIR := target
RELEASE_DIR := $(TARGET_DIR)/release
DEBUG_DIR := $(TARGET_DIR)/debug

.PHONY: gen build run clean test coverage

all: build test

$(DEBUG_DIR) $(RELEASE_DIR):
	@mkdir -p $@

compile: gen
	@go build ./pkg/...

gen:
	@go generate ./...

build: $(DEBUG_DIR) $(RELEASE_DIR) gen
	@for dir in $(shell find ./cmd -mindepth 1 -maxdepth 1 -type d); do \
		app_name=$$(basename $$dir); \
		go build --tags "release build" -o target/release/$$app_name $$dir; \
		go build --tags "debug build" -o target/debug/$$app_name $$dir; \
	done

test: $(DEBUG_DIR) $(RELEASE_DIR)
	@go test ./... -coverprofile=target/coverage.out

coverage: test
	@go tool cover -html=target/coverage.out -o=target/coverage.html

run:
	@go run ./cmd

clean:
	rm -rf $(TARGET_DIR)
