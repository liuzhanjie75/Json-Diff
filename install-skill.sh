#!/usr/bin/env sh
set -eu

PATH="/usr/bin:/bin:$PATH"
export PATH

case "$0" in
    */*) SCRIPT_BASE=${0%/*} ;;
    *) SCRIPT_BASE=. ;;
esac
SCRIPT_DIR=$(CDPATH= cd -- "$SCRIPT_BASE" && pwd)
cd "$SCRIPT_DIR"

GO_CMD=
if command -v go >/dev/null 2>&1; then
    GO_CMD=go
elif [ -x "/c/Program Files/Go/bin/go.exe" ]; then
    GO_CMD="/c/Program Files/Go/bin/go.exe"
fi

if [ -z "$GO_CMD" ]; then
    echo "Error: Go was not found in PATH." >&2
    exit 1
fi

CODEX_ROOT=${CODEX_HOME:-"$HOME/.codex"}
SKILLS_DIR="$CODEX_ROOT/skills"
TARGET_DIR="$SKILLS_DIR/jsondiff"
STAGING_DIR="$SKILLS_DIR/.jsondiff-install-$$"
BACKUP_DIR="$SKILLS_DIR/.jsondiff-backup-$$"
GOCACHE_DIR="$SCRIPT_DIR/.gocache-skill-install"
export GOCACHE="$GOCACHE_DIR"

cleanup() {
    if [ -e "$BACKUP_DIR" ] && [ ! -e "$TARGET_DIR" ]; then
        mv "$BACKUP_DIR" "$TARGET_DIR"
    fi
    rm -rf "$STAGING_DIR"
    rm -rf "$GOCACHE_DIR"
}

on_signal() {
    cleanup
    exit 1
}

trap cleanup EXIT
trap on_signal HUP INT TERM

mkdir -p "$STAGING_DIR/bin"
cp -R skill/jsondiff/. "$STAGING_DIR/"

OUTPUT="$STAGING_DIR/bin/jsondiff"
if [ "$("$GO_CMD" env GOOS)" = "windows" ]; then
    OUTPUT="$STAGING_DIR/bin/jsondiff.exe"
fi

echo "Building jsondiff for the skill..."
"$GO_CMD" build -o "$OUTPUT" .

if [ -e "$TARGET_DIR" ]; then
    mv "$TARGET_DIR" "$BACKUP_DIR"
fi

if ! mv "$STAGING_DIR" "$TARGET_DIR"; then
    if [ -e "$BACKUP_DIR" ]; then
        mv "$BACKUP_DIR" "$TARGET_DIR"
    fi
    exit 1
fi

rm -rf "$BACKUP_DIR"
rm -rf "$GOCACHE_DIR"
trap - EXIT HUP INT TERM
echo "Installed JSON Diff skill: $TARGET_DIR"
echo "Restart Codex or start a new session to discover the skill."
