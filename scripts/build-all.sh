#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="$ROOT_DIR/builds"
APP_NAME="pteryx"

VERSION="${PTERYX_VERSION:-}"
if [[ -z "$VERSION" ]]; then
  VERSION="$(git -C "$ROOT_DIR" describe --tags --abbrev=0 2>/dev/null || echo "dev")"
fi
VERSION="${VERSION#v}"

mkdir -p "$OUT_DIR"

build() {
  local goos="$1"
  local goarch="$2"
  local ext="${3:-}"
  local output="$OUT_DIR/${APP_NAME}-${VERSION}-${goos}-${goarch}${ext}"

  echo "Building ${goos}/${goarch} -> ${output}"
  (
    cd "$ROOT_DIR"
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
      go build -trimpath -ldflags="-s -w" -o "$output" .
  )
}

build darwin amd64
build darwin arm64
build linux amd64
build windows amd64 ".exe"

echo
echo "Built release artifacts in $OUT_DIR:"
find "$OUT_DIR" -maxdepth 1 -type f -name "${APP_NAME}-${VERSION}-*" -print | sort
