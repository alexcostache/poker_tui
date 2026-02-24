#!/usr/bin/env bash
set -euo pipefail

BINARY="poker_tui"
CMD="./cmd/poker_tui"
TARGET="/usr/local/bin/$BINARY"
ALIAS="/usr/local/bin/ptui"

# Build if missing
if [ ! -f "$BINARY" ]; then
    echo "==> Building $BINARY..."
    go build -o "$BINARY" "$CMD"
fi

# Copy to /usr/local/bin
sudo cp "$BINARY" "$TARGET"

# Create alias symlink
if [ ! -e "$ALIAS" ]; then
    sudo ln -s "$TARGET" "$ALIAS"
    echo "==> Symlinked $ALIAS to $TARGET"
fi

echo "==> Installed: $TARGET"
echo "==> You can now run 'poker_tui' or 'ptui' from any terminal."
