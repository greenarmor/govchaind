#!/bin/bash
# This script verifies that the local environment is correctly configured to build the project.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Verifying build environment for OpenGovChain...${NC}"

# --- Helper Functions ---
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# --- Dependency Checks ---

# 1. Check for Go
if ! command_exists go; then
    echo -e "${RED}❌ Go is not installed. Please install Go and ensure it is in your PATH.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Go is installed.${NC}"

# 2. Check Go version
REQUIRED_GO_VERSION=$(grep -E '^go [0-9.]+' go.mod | cut -d ' ' -f 2)
INSTALLED_GO_VERSION=$(go version | { read -r _ _ v _; echo "${v#go}"; })

if ! printf '%s\n' "$REQUIRED_GO_VERSION" "$INSTALLED_GO_VERSION" | sort -V -C; then
    echo -e "${RED}❌ Go version mismatch.${NC}"
    echo "  Required version: $REQUIRED_GO_VERSION"
    echo "  Installed version: $INSTALLED_GO_VERSION"
    echo "  Please upgrade or downgrade your Go installation."
    exit 1
fi
echo -e "${GREEN}✅ Go version is compatible ($INSTALLED_GO_VERSION).${NC}"

# 3. Check for Ignite CLI
if ! command_exists ignite; then
    echo -e "${RED}❌ Ignite CLI is not installed.${NC}"
    echo "  Please install it by following the official guide: https://docs.ignite.com/guide/install"
    exit 1
fi
echo -e "${GREEN}✅ Ignite CLI is installed.${NC}"

# 4. Check for C Compiler (gcc)
if ! command_exists gcc; then
    echo -e "${RED}❌ C compiler (gcc) is not installed.${NC}"
    echo "  On Debian/Ubuntu, run: sudo apt-get install build-essential"
    echo "  On macOS, run: xcode-select --install"
    exit 1
fi
echo -e "${GREEN}✅ C compiler (gcc) is installed.${NC}"

# 5. Check Go modules
echo "Verifying Go modules..."
go mod verify
echo -e "${GREEN}✅ Go modules are verified.${NC}"

# --- Success ---

echo ""
echo -e "${GREEN}🎉 Environment check passed! You are ready to build.${NC}"
_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
echo "Run 'ignite chain build' from the project root directory ('${_SCRIPT_DIR}/..') to build the binary."
