package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/srivastavya/goober/internal"
)

// TUI holds the Bubble Tea program and integrates with the runner
type TUI struct {
	program *tea.Program
	model   Model
	runner  *internal.Runner
}

// Init implements tea.Model
func (t TUI) Init() tea.Cmd {
	return t.model.Init()
}

// Update implements tea.Model
func (t TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle force restart at the TUI level
	if _, ok := msg.(ForceRestartMsg); ok {
		t.ForceRestart()
		return t, nil
	}
	
	// Delegate to the model
	updatedModel, cmd := t.model.Update(msg)
	t.model = updatedModel.(Model)
	return t, cmd
}

// View implements tea.Model
func (t TUI) View() string {
	return t.model.View()
}

// NewTUI creates a new TUI instance
func NewTUI(watchDir string, debounceDuration time.Duration, runner *internal.Runner) *TUI {
	model := NewModel(watchDir, debounceDuration)
	
	tui := &TUI{
		model:  model,
		runner: runner,
	}
	
	// Set up callbacks
	runner.SetLogCallback(tui.handleLog)
	runner.SetRestartCallback(tui.handleRestart)
	
	// Create the Bubble Tea program with a custom update function that handles force restart
	tui.program = tea.NewProgram(tui, tea.WithAltScreen(), tea.WithMouseCellMotion())
	
	return tui
}

// Start runs the TUI
func (t *TUI) Start() error {
	_, err := t.program.Run()
	return err
}

// Stop stops the TUI and the runner
func (t *TUI) Stop() {
	t.runner.Stop()
	t.program.Quit()
}

// handleLog processes log messages from the runner
func (t *TUI) handleLog(message string, isError bool) {
	logType := LogTypeInfo
	if isError {
		logType = LogTypeError
	}
	
	// Check for specific patterns to determine log type
	if contains(message, []string{"successful", "success", "✓", "completed"}) {
		logType = LogTypeSuccess
	} else if contains(message, []string{"warning", "warn", "⚠"}) {
		logType = LogTypeWarning
	} else if contains(message, []string{"error", "failed", "fail", "✗", "panic"}) {
		logType = LogTypeError
	}
	
	// Send log message through channel
	select {
	case t.model.logChan <- LogMessage{Message: message, Type: logType}:
	default:
		// Channel is full, skip this message
	}
	
	// Update build status based on message content
	if contains(message, []string{"build successful", "build completed"}) {
		t.program.Send(BuildStatusMsg{Status: BuildStatusSuccess})
	} else if contains(message, []string{"build failed", "build error"}) {
		t.program.Send(BuildStatusMsg{Status: BuildStatusFailed})
	}
}

// handleRestart processes restart events
func (t *TUI) handleRestart() {
	t.program.Send(RestartCountMsg{})
}

// contains checks if the message contains any of the given substrings
func contains(message string, substrings []string) bool {
	for _, substr := range substrings {
		if len(message) >= len(substr) {
			for i := 0; i <= len(message)-len(substr); i++ {
				match := true
				for j := 0; j < len(substr); j++ {
					if message[i+j] != substr[j] && message[i+j] != substr[j]-32 { // case insensitive
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
		}
	}
	return false
}

// GetProgram returns the Bubble Tea program for external control
func (t *TUI) GetProgram() *tea.Program {
	return t.program
}

// ForceRestart triggers a manual restart
func (t *TUI) ForceRestart() {
	go func() {
		t.handleLog("Manual restart triggered", false)
		t.runner.Stop()
		t.handleRestart()
		if err := t.runner.Start(); err != nil {
			t.handleLog("Manual restart failed: "+err.Error(), true)
		} else {
			t.handleLog("Manual restart successful", false)
		}
	}()
}
