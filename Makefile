BIN_DIR := bin

.PHONY: gen build run clean

all: build

gen:
	@go generate ./...

build: gen
	@mkdir -p $(BIN_DIR)
	@for dir in $(shell find ./cmd -type d -mindepth 1 -maxdepth 1); do \
		app_name=$$(basename $$dir); \
		go build --tags release -o $(BIN_DIR)/$$app_name $$dir; \
	done

run:
	@go run ./cmd

clean:
	@rm -rf $(BIN_DIR)
