package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jerkeyray/goober/internal"
	"github.com/jerkeyray/goober/tui"
)

func main() {
    dir := flag.String("dir", ".", "Directory to watch")
    buildCmd := flag.String("build", "go build -o app", "Build command")
    runCmd := flag.String("run", "./app", "Run command")
    debounce := flag.Duration("debounce", 750*time.Millisecond, "Debounce duration for file changes")
    noTUI := flag.Bool("no-tui", false, "Disable TUI and use simple CLI output")
    flag.Parse()

    runner := internal.NewRunner(*buildCmd, *runCmd)

    if *noTUI {
        // Original CLI mode
        fmt.Println("Starting Goober...")
        go runner.WatchAndRun(*dir)

        // Wait for Ctrl+C
        sig := make(chan os.Signal, 1)
        signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
        <-sig
        runner.Stop()
        fmt.Println("Goober stopped.")
    } else {
        // TUI mode
        tuiApp := tui.NewTUI(*dir, *debounce, runner)
        
        // Start the watcher in a goroutine
        go runner.WatchAndRun(*dir)
        
        // Start the TUI (this blocks until quit)
        if err := tuiApp.Start(); err != nil {
            fmt.Printf("Error running TUI: %v\n", err)
            os.Exit(1)
        }
        
        // Clean shutdown
        tuiApp.Stop()
    }
}
