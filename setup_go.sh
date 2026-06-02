#!/bin/bash
set -e

echo "=== Installing Go 1.23.4 ==="
mkdir -p "$HOME/go_dist"
curl -fsSL https://go.dev/dl/go1.23.4.linux-amd64.tar.gz -o /tmp/go.tar.gz
rm -rf "$HOME/go_dist/go"
tar -C "$HOME/go_dist" -xzf /tmp/go.tar.gz
rm /tmp/go.tar.gz

# Add Go to PATH in .bashrc if not already there
export PATH=$HOME/go_dist/go/bin:$PATH:$HOME/go/bin
if ! grep -q 'go_dist/go/bin' ~/.bashrc; then
    echo 'export PATH=$HOME/go_dist/go/bin:$PATH:$HOME/go/bin' >> ~/.bashrc
fi

echo "Go version: $($HOME/go_dist/go/bin/go version)"

# Install sqlc
echo "=== Installing sqlc ==="
$HOME/go_dist/go/bin/go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
echo "sqlc installed"

echo "=== Setup complete ==="
