#!/usr/bin/env bash
set -euo pipefail

BINARY="poker_tui"
CMD="./cmd/poker_tui"

echo "==> Building $BINARY..."
go build -o "$BINARY" "$CMD"
echo "==> Done: ./$BINARY"
