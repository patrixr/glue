#!/bin/bash

# Define variables
REPO="patrixr/glue"
API_URL="https://api.github.com/repos/$REPO/releases/latest"
TMP_DIR=$(mktemp -d)

# Fetch the latest release information
RELEASE_INFO=$(curl -s $API_URL)

# Determine the appropriate binary for the OS architecture
case "$(uname -s)" in
    Darwin)
        ARCH=$(uname -m)
        if [[ "$ARCH" == "x86_64" ]]; then
            BINARY_URL=$(echo "$RELEASE_INFO" | jq -r '.assets[] | select(.name | contains("darwin-amd64")) | .browser_download_url')
        elif [[ "$ARCH" == "arm64" ]]; then
            BINARY_URL=$(echo "$RELEASE_INFO" | jq -r '.assets[] | select(.name | contains("darwin-arm64")) | .browser_download_url')
        fi
        ;;
    Linux)
        ARCH=$(uname -m)
        if [[ "$ARCH" == "x86_64" ]]; then
            BINARY_URL=$(echo "$RELEASE_INFO" | jq -r '.assets[] | select(.name | contains("linux-amd64")) | .browser_download_url')
        elif [[ "$ARCH" == "arm64" ]]; then
            BINARY_URL=$(echo "$RELEASE_INFO" | jq -r '.assets[] | select(.name | contains("linux-arm64")) | .browser_download_url')
        fi
        ;;
    *)
        echo "Unsupported OS"
        exit 1
        ;;
esac

# Download and extract the binary
if [ -n "$BINARY_URL" ]; then
    echo "Downloading binary from $BINARY_URL"
    curl -L -o "$TMP_DIR/binary.tar.gz" "$BINARY_URL"
    tar -xzf "$TMP_DIR/binary.tar.gz" -C "$TMP_DIR"
else
    echo "No suitable binary found for architecture $ARCH"
    exit 1
fi

# Assuming the binary extracted is named 'glue', run it
if [ -f "$TMP_DIR/glue" ]; then
    echo "Running the binary..."
    "$TMP_DIR/glue"
else
    echo "Binary not found in extracted files"
    exit 1
fi

# Clean up
echo "Cleaning up..."
rm -rf "$TMP_DIR"

echo "Done."
