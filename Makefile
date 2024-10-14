.PHONY: tidy dry-run example help build test clean install uninstall

BIN_NAME := glue
BUILD_FOLDER := $(CURDIR)/.out
EXAMPLE_DIR := examples/homebrew
PREFIX ?= /usr/local
INSTALL_PATH = $(PREFIX)/bin

build:
	go build -o ${BUILD_FOLDER}/${BIN_NAME} ./

dry-run:
	@echo "Select test folder to run:"
	@select d in ./examples/*/; do test -n "$$d" && go run ./ --dry-run -p "$$d"; break; echo ">>> Invalid Selection"; done

example:
	@echo "Select test folder to run:"
	@select d in ./examples/*/; do test -n "$$d" && go run ./ -p "$$d"; break; echo ">>> Invalid Selection"; done

document:
	@go run ./ document

help:
	go run ./ --help

tidy:
	go mod tidy

test:
	ENV=test go test -json -v  ./...  | go run  github.com/mfridman/tparse@latest -all

clean:
	find . -name "*~" -delete
	find . -name ".DS_Store" -delete
	find . -name "#*" -delete

install: check-root
	@echo "Creating symlink for $(BIN_NAME) in $(DESTDIR)$(INSTALL_PATH)"
	@mkdir -p $(DESTDIR)$(INSTALL_PATH)
	@ln -sf $(BUILD_FOLDER)/$(BIN_NAME) $(DESTDIR)$(INSTALL_PATH)/$(BIN_NAME)
	@echo "Installation complete. You can now run '$(BIN_NAME)'"

uninstall: check-root
	@echo "Removing symlink for $(BIN_NAME) from $(DESTDIR)$(INSTALL_PATH)"
	@rm -f $(DESTDIR)$(INSTALL_PATH)/$(BIN_NAME)
	@echo "Uninstallation complete"

check-root:
	@if [ $$(id -u) -ne 0 ]; then \
		echo "This target must be run as root or with sudo."; \
		exit 1; \
	fi
