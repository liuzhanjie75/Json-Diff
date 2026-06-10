#!/usr/bin/env sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
cd "$SCRIPT_DIR"

if ! command -v go >/dev/null 2>&1; then
    echo "Error: Go was not found in PATH." >&2
    exit 1
fi

OUTPUT=jsondiff
if [ "$(go env GOOS)" = "windows" ]; then
    OUTPUT=jsondiff.exe
fi

echo "Building $OUTPUT..."
go build -o "$OUTPUT" .
echo "Build complete: $SCRIPT_DIR/$OUTPUT"
