# Goober - Go Development Watcher with Beautiful TUI 🎯

A modern file watcher that automatically rebuilds and restarts your Go applications with a beautiful Bubble Tea terminal interface.

![Goober Demo](https://raw.githubusercontent.com/yourusername/goober/main/demo.gif)

## ✨ Features

- 🔄 **Auto-reload** - Watches `.go` files and restarts on changes
- 🎨 **Beautiful TUI** - Modern terminal interface with real-time logs
- 🌈 **Color-coded output** - Success, error, and warning messages
- 📜 **Scrollable logs** - Navigate through build and runtime output
- ⚡ **Manual restart** - Force restart with a keypress
- 🧹 **Log management** - Clear logs on demand
- 🎛️ **Configurable** - Customizable build/run commands and debounce timing

## 🚀 Quick Start

### Global Installation

```bash
go install github.com/srivastavya/goober/cmd/goober@latest
```

### Usage in Your Go Project

```bash
# Navigate to your Go project
cd your-go-project

# Start with default settings (builds to 'app', runs './app')
goober

# Custom build and run commands
goober -build "go build -o myserver ./cmd/server" -run "./myserver --port 8080"

# Watch specific directory
goober -dir ./src -debounce 1s
```

## 📖 Installation Options

### 1. Install Globally (Recommended)

```bash
go install github.com/srivastavya/goober/cmd/goober@latest
```

Then use `goober` anywhere in your terminal.

### 2. Install in Project

Add to your project's `tools.go`:

```go
//go:build tools

package tools

import _ "github.com/srivastavya/goober/cmd/goober"
```

Then install:

```bash
go mod tidy
go install github.com/srivastavya/goober/cmd/goober
```

### 3. Download Binary

Download from [releases page](https://github.com/yourusername/goober/releases) and add to PATH.

### 4. Build from Source

```bash
git clone https://github.com/srivastavya/goober.git
cd goober
go build -o goober ./cmd/goober
sudo mv goober /usr/local/bin/
```

## Usage

### TUI Mode (Default)

```bash
# Watch current directory with default settings
./goober

# Watch specific directory with custom commands
./goober -dir ./myapp -build "go build -o myapp" -run "./myapp"

# Custom debounce duration
./goober -debounce 1s
```

### CLI Mode

```bash
# Disable TUI for simple terminal output
./goober -no-tui
```

## TUI Interface

### Status Bar (Top)

- **Directory**: Currently watched directory
- **Debounce**: File change debounce duration
- **Restarts**: Number of automatic restarts
- **Build**: Last build status (Success/Failed/Unknown)

### Log Panel (Middle)

- Real-time build and application output
- Color-coded messages:
  - 🟢 **Green**: Success messages
  - 🔴 **Red**: Error messages
  - 🟡 **Yellow**: Warning messages
  - ⚪ **Gray**: Info messages
- Scrollable with timestamps
- Shows scroll position when logs exceed screen

### Help Bar (Bottom)

- **q / Ctrl+C**: Quit application
- **r**: Force manual restart
- **c**: Clear all logs
- **↑/↓ or j/k**: Scroll through logs
- **PgUp/PgDown**: Fast scroll

## Command Line Options

- `-dir string`: Directory to watch (default: current directory)
- `-build string`: Build command (default: "go build -o app")
- `-run string`: Run command (default: "./app")
- `-debounce duration`: Debounce duration for file changes (default: 750ms)
- `-no-tui`: Disable TUI and use simple CLI output

## Examples

### Basic Web Server

```bash
./goober -dir ./server -build "go build -o server ./cmd/server" -run "./server"
```

### With Custom Build Tags

```bash
./goober -build "go build -tags dev -o app" -run "./app --dev"
```

### Microservice Development

```bash
./goober -dir ./service -debounce 1s -build "go build -o service" -run "./service --port 8080"
```

## Architecture

The application is structured into several packages:

- `cmd/goober/`: Main application entry point
- `internal/`: Core functionality (runner and watcher)
- `tui/`: Bubble Tea terminal user interface
  - `model.go`: Application state and data structures
  - `view.go`: UI rendering and layout
  - `update.go`: Event handling and state updates
  - `tui.go`: Integration layer between TUI and backend

## Development

The TUI is built using:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling and layout

To contribute:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test with both TUI and CLI modes
5. Submit a pull request

## License

[Your License Here]
