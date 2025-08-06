#!/bin/bash
# Goober Installation Script
set -e

REPO="srivastavya/goober"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🎯 Goober Installation Script${NC}"
echo "================================"

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
  x86_64) ARCH="x86_64" ;;
  arm64|aarch64) ARCH="arm64" ;;
  *) 
    echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"
    exit 1
    ;;
esac

case $OS in
  linux) OS="Linux" ;;
  darwin) OS="Darwin" ;;
  *) 
    echo -e "${RED}❌ Unsupported OS: $OS${NC}"
    exit 1
    ;;
esac

echo -e "${YELLOW}📋 Detected: $OS/$ARCH${NC}"

# Get latest release
echo -e "${YELLOW}🔍 Fetching latest release...${NC}"
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep -o '"tag_name": "[^"]*"' | cut -d'"' -f4)

if [ -z "$LATEST_RELEASE" ]; then
  echo -e "${RED}❌ Failed to fetch latest release${NC}"
  echo -e "${YELLOW}💡 Try installing via Go: go install github.com/$REPO/cmd/goober@latest${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Latest version: $LATEST_RELEASE${NC}"

# Download URL
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_RELEASE/goober_${OS}_${ARCH}.tar.gz"

echo -e "${YELLOW}📥 Downloading from: $DOWNLOAD_URL${NC}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download and extract
if curl -L --fail "$DOWNLOAD_URL" | tar xz; then
  echo -e "${GREEN}✅ Download and extraction successful${NC}"
else
  echo -e "${RED}❌ Download failed${NC}"
  echo -e "${YELLOW}💡 Try installing via Go: go install github.com/$REPO/cmd/goober@latest${NC}"
  exit 1
fi

# Install
echo -e "${YELLOW}📦 Installing to $INSTALL_DIR...${NC}"

if [ -w "$INSTALL_DIR" ]; then
  mv goober "$INSTALL_DIR/"
  echo -e "${GREEN}✅ Goober installed successfully!${NC}"
else
  echo -e "${YELLOW}🔐 Requesting sudo access to install to $INSTALL_DIR...${NC}"
  sudo mv goober "$INSTALL_DIR/"
  echo -e "${GREEN}✅ Goober installed successfully with sudo!${NC}"
fi

# Cleanup
cd - > /dev/null
rm -rf "$TMP_DIR"

# Verify installation
if command -v goober >/dev/null 2>&1; then
  echo -e "${GREEN}🎉 Installation verified!${NC}"
  echo -e "${BLUE}📍 Goober location: $(which goober)${NC}"
  echo -e "${BLUE}📖 Version: $(goober -h 2>&1 | head -1 || echo 'Goober file watcher')${NC}"
else
  echo -e "${YELLOW}⚠️  Installation completed but 'goober' not found in PATH${NC}"
  echo -e "${YELLOW}💡 You may need to add $INSTALL_DIR to your PATH${NC}"
fi

echo ""
echo -e "${GREEN}🚀 Quick Start:${NC}"
echo -e "  ${BLUE}cd your-go-project${NC}"
echo -e "  ${BLUE}goober${NC}"
echo ""
echo -e "${GREEN}📖 More examples:${NC}"
echo -e "  ${BLUE}goober -build 'go build -o app' -run './app'${NC}"
echo -e "  ${BLUE}goober -dir ./src -debounce 1s${NC}"
echo -e "  ${BLUE}goober --help${NC}"
echo ""
echo -e "${GREEN}🎯 Happy coding with Goober!${NC}"
