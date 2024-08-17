.PHONY: tidy run build test clean

BIN_NAME := glue
BUILD_FOLDER := ./.out

build:
	go build -o ${BUILD_FOLDER}/${BIN_NAME} ./

run:
	go run ./

tidy:
	go mod tidy

test:
	ENV=test go test  -v ./...

clean:
	find . -name "*~" -delete
	find . -name ".DS_Store" -delete
	find . -name "#*" -delete
