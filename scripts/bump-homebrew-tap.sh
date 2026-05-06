#!/usr/bin/env bash
set -euo pipefail

cd "$(brew --repository dtsivkovski/homebrew-tap)"
scripts/bump-pteryx.sh "$@"