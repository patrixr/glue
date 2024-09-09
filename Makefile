.PHONY: tidy dry help build test clean

BIN_NAME := glue
BUILD_FOLDER := ./.out

build:
	go build -o ${BUILD_FOLDER}/${BIN_NAME} ./

dry:
	go run ./ --dry-run

example:
	go run ./ -p examples/unsafe

#go run ./ -f examples/prints/glue.lue

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
