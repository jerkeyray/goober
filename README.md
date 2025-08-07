# Goober

A badass file watcher for Go projects that automatically rebuilds and restarts your application. It comes with a slick terminal UI to show logs and build status.

## 🚀 Installation

Install Goober globally:

```bash
go install github.com/jerkeyray/goober@latest
```

Make sure your Go bin directory (usually `$HOME/go/bin`) is in your `PATH`. If not:

```bash
export PATH="$HOME/go/bin:$PATH"
```

Add that to your `.zshrc`, `.bashrc`, or whatever shell config you use.

## ⚡ Quick Start

Navigate to your Go project and just run:

```bash
goober
```

It will:

- Watch for file changes
- Rebuild the project
- Restart your app
- Show logs in a terminal UI

### Examples

```bash
# Watch current directory and use default build/run
goober

# Watch a specific directory
goober --dir ./myapp

# Use custom build and run commands
goober --build "go build -o myapp" --run "./myapp"

# Set a custom debounce time (e.g., 1 second)
goober --debounce 1s

# Disable the TUI, use plain logs
goober --no-tui
```

## 🛠️ CLI Flags

- `--dir <path>` — Directory to watch (default: `.`)
- `--build <command>` — Build command (default: `go build -o app`)
- `--run <command>` — Run command (default: `./app`)
- `--debounce <duration>` — Delay after file changes before restarting (default: `750ms`)
- `--no-tui` — Disable the terminal UI and use plain output

## 🎮 TUI Keybindings

When using the terminal UI:

- `q` / `Ctrl+C` — Quit
- `r` — Force manual restart
- `c` — Clear logs
- `↑`/`↓` or `j`/`k` — Scroll logs
- `PgUp` / `PgDown` — Scroll a page up/down

## 📄 License

MIT

---

Made with rage and caffeine by [@jerkeyray](https://github.com/jerkeyray)
