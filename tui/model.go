package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// BuildStatus represents the status of the last build
type BuildStatus int

const (
	BuildStatusUnknown BuildStatus = iota
	BuildStatusSuccess
	BuildStatusFailed
)

// LogEntry represents a single log entry with timestamp and styling
type LogEntry struct {
	Message   string
	Type      LogType
	Timestamp time.Time
}

// LogType determines styling for log entries
type LogType int

const (
	LogTypeInfo LogType = iota
	LogTypeSuccess
	LogTypeError
	LogTypeWarning
)

// Model represents the state of our TUI application
type Model struct {
	// Status bar data
	watchDir        string
	debounceDuration time.Duration
	restartCount    int
	buildStatus     BuildStatus
	
	// Logs
	logs         []LogEntry
	logViewStart int // For scrolling
	
	// UI state
	width         int
	height        int
	ready         bool
	
	// Styles
	styles Styles
	
	// Channel for receiving log messages
	logChan chan LogMessage
}

// LogMessage is sent through the channel to add new log entries
type LogMessage struct {
	Message string
	Type    LogType
}

// Styles holds all the styling for the TUI
type Styles struct {
	StatusBar       lipgloss.Style
	StatusBarKey    lipgloss.Style
	StatusBarValue  lipgloss.Style
	LogPanel        lipgloss.Style
	LogEntryInfo    lipgloss.Style
	LogEntrySuccess lipgloss.Style
	LogEntryError   lipgloss.Style
	LogEntryWarning lipgloss.Style
	HelpBar         lipgloss.Style
	HelpKey         lipgloss.Style
	HelpDesc        lipgloss.Style
}

// NewModel creates a new model with default values
func NewModel(watchDir string, debounceDuration time.Duration) Model {
	return Model{
		watchDir:        watchDir,
		debounceDuration: debounceDuration,
		restartCount:    0,
		buildStatus:     BuildStatusUnknown,
		logs:            []LogEntry{},
		logViewStart:    0,
		width:          80,
		height:         24,
		ready:          false,
		styles:         makeStyles(),
		logChan:        make(chan LogMessage, 100),
	}
}

// makeStyles creates the styling configuration
func makeStyles() Styles {
	return Styles{
		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#7C3AED")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		
		StatusBarKey: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A78BFA")).
			Bold(true),
		
		StatusBarValue: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		
		LogPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2),
		
		LogEntryInfo: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6B7280")),
		
		LogEntrySuccess: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")),
		
		LogEntryError: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")),
		
		LogEntryWarning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F59E0B")),
		
		HelpBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#374151")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),
		
		HelpKey: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60A5FA")).
			Bold(true),
		
		HelpDesc: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D1D5DB")),
	}
}

// AddLog adds a new log entry to the model
func (m *Model) AddLog(message string, logType LogType) {
	entry := LogEntry{
		Message:   message,
		Type:      logType,
		Timestamp: time.Now(),
	}
	m.logs = append(m.logs, entry)
	
	// Auto-scroll to bottom when new logs are added
	maxVisible := m.getMaxVisibleLogs()
	if len(m.logs) > maxVisible {
		m.logViewStart = len(m.logs) - maxVisible
	}
}

// ClearLogs removes all log entries
func (m *Model) ClearLogs() {
	m.logs = []LogEntry{}
	m.logViewStart = 0
}

// IncrementRestartCount increases the restart counter
func (m *Model) IncrementRestartCount() {
	m.restartCount++
}

// SetBuildStatus updates the build status
func (m *Model) SetBuildStatus(status BuildStatus) {
	m.buildStatus = status
}

// GetLogChan returns the channel for sending log messages
func (m *Model) GetLogChan() chan LogMessage {
	return m.logChan
}

// getMaxVisibleLogs calculates how many log lines can fit in the current view
func (m *Model) getMaxVisibleLogs() int {
	// Height minus status bar, help bar, and log panel borders/padding
	return m.height - 6
}

// getVisibleLogs returns the logs that should be displayed in the current view
func (m *Model) getVisibleLogs() []LogEntry {
	maxVisible := m.getMaxVisibleLogs()
	totalLogs := len(m.logs)
	
	if totalLogs == 0 {
		return []LogEntry{}
	}
	
	start := m.logViewStart
	end := start + maxVisible
	
	if end > totalLogs {
		end = totalLogs
	}
	
	return m.logs[start:end]
}

// scrollUp moves the log view up
func (m *Model) scrollUp() {
	if m.logViewStart > 0 {
		m.logViewStart--
	}
}

// scrollDown moves the log view down
func (m *Model) scrollDown() {
	maxVisible := m.getMaxVisibleLogs()
	maxStart := len(m.logs) - maxVisible
	if maxStart < 0 {
		maxStart = 0
	}
	
	if m.logViewStart < maxStart {
		m.logViewStart++
	}
}

// formatLogEntry formats a log entry for display
func (m *Model) formatLogEntry(entry LogEntry) string {
	timestamp := entry.Timestamp.Format("15:04:05")
	
	var style lipgloss.Style
	switch entry.Type {
	case LogTypeSuccess:
		style = m.styles.LogEntrySuccess
	case LogTypeError:
		style = m.styles.LogEntryError
	case LogTypeWarning:
		style = m.styles.LogEntryWarning
	default:
		style = m.styles.LogEntryInfo
	}
	
	return fmt.Sprintf("%s %s", 
		lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF")).Render(timestamp),
		style.Render(entry.Message))
}

// buildStatusString returns a formatted string for the build status
func (m *Model) buildStatusString() string {
	switch m.buildStatus {
	case BuildStatusSuccess:
		return m.styles.LogEntrySuccess.Render("Success")
	case BuildStatusFailed:
		return m.styles.LogEntryError.Render("Failed")
	default:
		return m.styles.StatusBarValue.Render("Unknown")
	}
}
