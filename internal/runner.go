package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type LogCallback func(message string, isError bool)
type RestartCallback func()

type Runner struct {
	buildCmd        string
	runCmd          string
	proc            *exec.Cmd
	mu              sync.Mutex
	logCallback     LogCallback
	restartCallback RestartCallback
}

func NewRunner(build, run string) *Runner {
	return &Runner{buildCmd: build, runCmd: run}
}

func (r *Runner) SetLogCallback(callback LogCallback) {
	r.logCallback = callback
}

func (r *Runner) SetRestartCallback(callback RestartCallback) {
	r.restartCallback = callback
}

func (r *Runner) log(message string, isError bool) {
	if r.logCallback != nil {
		r.logCallback(message, isError)
	} else {
		if isError {
			fmt.Fprintf(os.Stderr, "%s\n", message)
		} else {
			fmt.Println(message)
		}
	}
}

// Build and run the app
func (r *Runner) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.log("Building...", false)
	if err := r.runCommand(r.buildCmd); err != nil {
		r.log(fmt.Sprintf("Build failed: %v", err), true)
		return err
	}

	r.log("Build successful", false)
	r.log("Starting app...", false)
	
	parts := strings.Split(r.runCmd, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	
	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	
	if err := cmd.Start(); err != nil {
		r.log(fmt.Sprintf("Failed to start app: %v", err), true)
		return err
	}
	
	r.proc = cmd
	
	// Stream output
	go r.streamOutput(stdout, false)
	go r.streamOutput(stderr, true)
	
	return nil
}

// streamOutput reads from a pipe and sends output to the log callback
func (r *Runner) streamOutput(pipe io.ReadCloser, isError bool) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			r.log(line, isError)
		}
	}
	pipe.Close()
}

// Stop the running app
func (r *Runner) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.proc == nil || r.proc.Process == nil {
		return
	}

	r.log("Stopping app...", false)
	if err := r.proc.Process.Signal(os.Interrupt); err != nil {
		r.log("Graceful stop failed, killing...", true)
		r.proc.Process.Kill()
	}
	r.proc.Wait()
	r.proc = nil
}

// Utility to run build commands with output capture
func (r *Runner) runCommand(cmdStr string) error {
	parts := strings.Split(cmdStr, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	// Log stdout
	if stdout.Len() > 0 {
		r.log(strings.TrimSpace(stdout.String()), false)
	}
	
	// Log stderr
	if stderr.Len() > 0 {
		r.log(strings.TrimSpace(stderr.String()), true)
	}
	
	return err
}
