.PHONY: tidy dry help build test clean

BIN_NAME := glue
BUILD_FOLDER := ./.out
BUILD_CMD :=

build: build-macos

build-macos:
	env GOOS=darwin GOARCH=amd64 go build -o ${BUILD_FOLDER}/macos-amd64/${BIN_NAME} ./
	env GOOS=darwin GOARCH=arm64 go build -o ${BUILD_FOLDER}/macos-arm64/${BIN_NAME} ./

dry:
	go run ./ --dry-run

example:
	go run ./ -p examples/unsafe

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
