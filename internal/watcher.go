package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (r *Runner) WatchAndRun(dir string) {
	if err := r.Start(); err != nil {
		r.log(fmt.Sprintf("Initial start failed: %v", err), true)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		r.log(fmt.Sprintf("Failed to create watcher: %v", err), true)
		return
	}
	defer watcher.Close()

	// Add directories recursively
	addDirs(watcher, dir)
	r.log(fmt.Sprintf("Watching directory: %s", dir), false)

	debounce := time.NewTimer(time.Hour)
	debounce.Stop()
	mu := sync.Mutex{}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if strings.HasSuffix(event.Name, ".go") {
					r.log(fmt.Sprintf("Change detected: %s", event.Name), false)
					mu.Lock()
					debounce.Reset(750 * time.Millisecond)
					mu.Unlock()
				}
			case err, ok := <-watcher.Errors:
				if ok {
					r.log(fmt.Sprintf("Watcher error: %v", err), true)
				}
			}
		}
	}()

	for {
		<-debounce.C
		r.log("Changes detected, restarting...", false)
		r.Stop()
		
		// Notify about restart
		if r.restartCallback != nil {
			r.restartCallback()
		}
		
		if err := r.Start(); err != nil {
			r.log(fmt.Sprintf("Restart failed: %v", err), true)
		} else {
			r.log("Restart successful", false)
		}
	}
}

func addDirs(watcher *fsnotify.Watcher, root string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && !strings.Contains(path, ".git") && !strings.Contains(path, "vendor") {
			watcher.Add(path)
		}
		return nil
	})
}
