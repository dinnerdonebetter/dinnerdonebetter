#!/bin/bash

ensure_wire_installed() {
    if ! command -v wire >/dev/null 2>&1; then
        echo "Installing wire..."
        go install github.com/google/wire/cmd/wire@v0.6.0
    fi
}
