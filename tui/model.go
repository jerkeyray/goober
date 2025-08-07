package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// BuildStatus represents the status of the last build
type BuildStatus int

const (
	StatusWatching BuildStatus = iota
	StatusBuilding
	StatusSuccess
	StatusError
)

// LogType determines styling for log entries
type LogType int

const (
	LogTypeInfo LogType = iota
	LogTypeSuccess
	LogTypeError
	LogTypeRestart
	LogTypeEvent
)

// LogEntry represents a single log entry with timestamp and styling
type LogEntry struct {
	Message   string
	Type      LogType
	Timestamp time.Time
}

// Model represents the state of our TUI application
type Model struct {
	// Status bar data
	watchDir         string
	debounceDuration time.Duration
	restartCount     int
	status           BuildStatus

	// Logs
	logs         []LogEntry
	logViewStart int // For scrolling

	// UI state
	width  int
	height int
	ready  bool

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
	StatusWatching  lipgloss.Style
	StatusBuilding  lipgloss.Style
	StatusSuccess   lipgloss.Style
	StatusError     lipgloss.Style
	LogPanel        lipgloss.Style
	LogTimestamp    lipgloss.Style
	LogEntryInfo    lipgloss.Style
	LogEntrySuccess lipgloss.Style
	LogEntryError   lipgloss.Style
	LogEntryRestart lipgloss.Style
	LogEntryEvent   lipgloss.Style
	HelpBar         lipgloss.Style
	HelpKey         lipgloss.Style
	HelpDesc        lipgloss.Style
}

// NewModel creates a new model with default values
func NewModel(watchDir string, debounceDuration time.Duration) Model {
	return Model{
		watchDir:         watchDir,
		debounceDuration: debounceDuration,
		restartCount:     0,
		status:           StatusWatching,
		logs:             []LogEntry{},
		logViewStart:     0,
		width:            80,
		height:           24,
		ready:            false,
		styles:           makeStyles(),
		logChan:          make(chan LogMessage, 100),
	}
}

// makeStyles creates the styling configuration based on the Tokyo Night theme
func makeStyles() Styles {
	// Tokyo Night Color Palette
	bg := lipgloss.Color("#1a1b26")
	fg := lipgloss.Color("#c0caf5")
	darkFg := lipgloss.Color("#a9b1d6")
	comment := lipgloss.Color("#565f89")
	blue := lipgloss.Color("#7aa2f7")
	purple := lipgloss.Color("#bb9af7")
	green := lipgloss.Color("#9ece6a")
	red := lipgloss.Color("#f7768e")
	yellow := lipgloss.Color("#e0af68")
	orange := lipgloss.Color("#ff9e64")

	return Styles{
		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#1f2335")).
			Foreground(darkFg).
			Padding(0, 1),
		StatusBarKey:   lipgloss.NewStyle().Foreground(purple).Bold(true),
		StatusBarValue: lipgloss.NewStyle().Foreground(fg),
		StatusWatching: lipgloss.NewStyle().Foreground(blue),
		StatusBuilding: lipgloss.NewStyle().Foreground(yellow),
		StatusSuccess:  lipgloss.NewStyle().Foreground(green),
		StatusError:    lipgloss.NewStyle().Foreground(red),

		LogPanel: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(1, 2),
		LogTimestamp:    lipgloss.NewStyle().Foreground(comment),
		LogEntryInfo:    lipgloss.NewStyle().Foreground(darkFg),
		LogEntrySuccess: lipgloss.NewStyle().Foreground(green),
		LogEntryError:   lipgloss.NewStyle().Foreground(red),
		LogEntryRestart: lipgloss.NewStyle().Foreground(orange),
		LogEntryEvent:   lipgloss.NewStyle().Foreground(blue),

		HelpBar:  lipgloss.NewStyle().Background(bg).Foreground(darkFg).Padding(0, 1),
		HelpKey:  lipgloss.NewStyle().Foreground(blue).Bold(true),
		HelpDesc: lipgloss.NewStyle().Foreground(comment),
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
	m.status = status
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

// formatLogEntry styles a single log entry
func (m Model) formatLogEntry(entry LogEntry) string {
	timestamp := m.styles.LogTimestamp.Render(entry.Timestamp.Format("15:04:05.000"))
	var logType lipgloss.Style
	var typeString string

	switch entry.Type {
	case LogTypeInfo:
		logType = m.styles.LogEntryInfo
		typeString = "INFO"
	case LogTypeSuccess:
		logType = m.styles.LogEntrySuccess
		typeString = "SUCCESS"
	case LogTypeError:
		logType = m.styles.LogEntryError
		typeString = "ERROR"
	case LogTypeRestart:
		logType = m.styles.LogEntryRestart
		typeString = "RESTART"
	case LogTypeEvent:
		logType = m.styles.LogEntryEvent
		typeString = "EVENT"
	default:
		logType = m.styles.LogEntryInfo
		typeString = "LOG"
	}

	typeLabel := logType.Copy().Bold(true).Render(fmt.Sprintf("[%s]", typeString))
	message := logType.Render(entry.Message)

	return fmt.Sprintf("%s %s %s", timestamp, typeLabel, message)
}

// buildStatusString returns a styled string for the current build status
func (m Model) buildStatusString() string {
	switch m.status {
	case StatusWatching:
		return m.styles.StatusWatching.Render("Watching")
	case StatusBuilding:
		return m.styles.StatusBuilding.Render("Building...")
	case StatusSuccess:
		return m.styles.StatusSuccess.Render("Running")
	case StatusError:
		return m.styles.StatusError.Render("Build Failed")
	default:
		return "Unknown"
	}
}
