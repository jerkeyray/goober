#!/bin/bash

# Demo script for Goober TUI
echo "🎯 Starting Goober TUI Demo"
echo "================================"
echo ""
echo "This demo will:"
echo "1. Start Goober with TUI watching the test/ directory"
echo "2. Build and run the demo web server"
echo "3. Show real-time logs in a beautiful interface"
echo ""
echo "TUI Controls:"
echo "  q / Ctrl+C  → Quit"
echo "  r           → Force restart"
echo "  c           → Clear logs"
echo "  ↑/↓         → Scroll logs"
echo ""
echo "Try editing test/main.go to see automatic reloading!"
echo ""

# Navigate to project root
cd "$(dirname "$0")"

# Run Goober with test directory
./goober -dir ./test -build "go build -o test/app test/main.go" -run "./test/app"
