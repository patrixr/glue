.PHONY: tidy dry help build test clean

BIN_NAME := glue
BUILD_FOLDER := ./.out
EXAMPLE_DIR := examples/homebrew

build:
	go build -o ${BUILD_FOLDER}/${BIN_NAME} ./

dry-run:
	@echo "Select test folder to run:"
	@select d in ./examples/*/; do test -n "$$d" && go run ./ --dry-run -p "$$d"; break; echo ">>> Invalid Selection"; done

example:
	@echo "Select test folder to run:"
	@select d in ./examples/*/; do test -n "$$d" && go run ./ -p "$$d"; break; echo ">>> Invalid Selection"; done

help:
	go run ./ --help

tidy:
	go mod tidy

test:
	ENV=test go test  -v ./...

clean:
	find . -name "*~" -delete
	find . -name ".DS_Store" -delete
	find . -name "#*" -delete
