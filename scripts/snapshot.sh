#!/usr/bin/env bash
#
# snapshot.sh - Take a snapshot of cannabis data using dank-extract
#
# Usage: ./scripts/snapshot.sh [YYYY-MM-DD]
#
# If no date is provided, uses today's date.
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Snapshot date (default: today)
SNAPSHOT_DATE="${1:-$(date +%Y-%m-%d)}"

echo "=== dank-data snapshot ==="
echo "Date: $SNAPSHOT_DATE"
echo "Repo: $REPO_ROOT"
echo ""

# Download latest dank-extract if not present or if --download flag
download_dank_extract() {
    echo "Downloading latest dank-extract..."

    local os arch download_url
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    arch="$(uname -m)"

    # Map architecture names
    case "$arch" in
        x86_64) arch="x86_64" ;;
        amd64)  arch="x86_64" ;;
        arm64)  arch="arm64" ;;
        aarch64) arch="arm64" ;;
    esac

    # Get download URL from latest release
    download_url=$(curl -s https://api.github.com/repos/AgentDank/dank-extract/releases/latest \
        | grep "browser_download_url" \
        | grep "${os}_${arch}" \
        | cut -d '"' -f 4)

    if [ -z "$download_url" ]; then
        echo "Error: Could not find dank-extract release for ${os}_${arch}"
        exit 1
    fi

    echo "Downloading from: $download_url"
    curl -L -o /tmp/dank-extract.tar.gz "$download_url"

    # Extract
    cd /tmp
    tar -xzf dank-extract.tar.gz
    mv */dank-extract "$REPO_ROOT/dank-extract"
    chmod +x "$REPO_ROOT/dank-extract"
    rm -rf dank-extract.tar.gz dank-extract_*/
    cd "$REPO_ROOT"

    echo "Downloaded dank-extract to $REPO_ROOT/dank-extract"
}

# Check if dank-extract exists, download if not
if [ ! -x "$REPO_ROOT/dank-extract" ]; then
    download_dank_extract
fi

# Run snapshot
echo ""
echo "Running snapshot..."
"$REPO_ROOT/dank-extract" --snapshot "$REPO_ROOT/snapshots" --snapshot-date "$SNAPSHOT_DATE" --verbose

# Update latest symlink
echo ""
echo "Updating latest symlink..."
cd "$REPO_ROOT/snapshots/us/ct"
rm -f latest
ln -s "$SNAPSHOT_DATE" latest
cd "$REPO_ROOT"

echo ""
echo "=== Snapshot complete ==="
echo "Files created in: $REPO_ROOT/snapshots/us/ct/$SNAPSHOT_DATE/"
ls -la "$REPO_ROOT/snapshots/us/ct/$SNAPSHOT_DATE/"
