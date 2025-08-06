package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the entire TUI
func (m Model) View() string {
	if !m.ready {
		return "Initializing Goober TUI..."
	}

	statusBar := m.renderStatusBar()
	helpBar := m.renderHelpBar()
	logPanel := m.renderLogPanel()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		statusBar,
		logPanel,
		helpBar,
	)
}

// renderStatusBar creates the top status bar
func (m Model) renderStatusBar() string {
	watchDirItem := fmt.Sprintf("%s %s", 
		m.styles.StatusBarKey.Render("Directory:"), 
		m.styles.StatusBarValue.Render(m.watchDir))

	debounceItem := fmt.Sprintf("%s %s", 
		m.styles.StatusBarKey.Render("Debounce:"), 
		m.styles.StatusBarValue.Render(m.debounceDuration.String()))

	restartItem := fmt.Sprintf("%s %s", 
		m.styles.StatusBarKey.Render("Restarts:"), 
		m.styles.StatusBarValue.Render(fmt.Sprintf("%d", m.restartCount)))

	buildStatusItem := fmt.Sprintf("%s %s", 
		m.styles.StatusBarKey.Render("Build:"), 
		m.buildStatusString())

	// Join status items with separators
	statusContent := strings.Join([]string{
		watchDirItem,
		debounceItem,
		restartItem,
		buildStatusItem,
	}, " │ ")

	// Apply styling and fit to width
	statusBar := m.styles.StatusBar.Render(statusContent)
	
	// Ensure the status bar spans the full width
	statusBar = m.styles.StatusBar.Width(m.width - 2).Render(statusContent)
	
	return statusBar
}

// renderLogPanel creates the scrollable log panel
func (m Model) renderLogPanel() string {
	visibleLogs := m.getVisibleLogs()
	
	if len(visibleLogs) == 0 {
		emptyMessage := m.styles.LogEntryInfo.Render("No logs yet. Watching for changes...")
		return m.styles.LogPanel.
			Width(m.width - 4).
			Height(m.height - 6).
			Render(emptyMessage)
	}

	var logLines []string
	for _, entry := range visibleLogs {
		logLines = append(logLines, m.formatLogEntry(entry))
	}

	// Add scroll indicator if there are more logs
	totalLogs := len(m.logs)
	maxVisible := m.getMaxVisibleLogs()
	
	var scrollInfo string
	if totalLogs > maxVisible {
		start := m.logViewStart + 1
		end := m.logViewStart + len(visibleLogs)
		scrollInfo = fmt.Sprintf("Showing %d-%d of %d logs", start, end, totalLogs)
		scrollIndicator := m.styles.LogEntryInfo.Render(scrollInfo)
		logLines = append(logLines, "", scrollIndicator)
	}

	logContent := strings.Join(logLines, "\n")

	return m.styles.LogPanel.
		Width(m.width - 4).
		Height(m.height - 6).
		Render(logContent)
}

// renderHelpBar creates the bottom help bar
func (m Model) renderHelpBar() string {
	helpItems := []string{
		fmt.Sprintf("%s %s", 
			m.styles.HelpKey.Render("q/Ctrl+C"), 
			m.styles.HelpDesc.Render("Quit")),
		fmt.Sprintf("%s %s", 
			m.styles.HelpKey.Render("r"), 
			m.styles.HelpDesc.Render("Restart")),
		fmt.Sprintf("%s %s", 
			m.styles.HelpKey.Render("c"), 
			m.styles.HelpDesc.Render("Clear logs")),
		fmt.Sprintf("%s %s", 
			m.styles.HelpKey.Render("↑/↓"), 
			m.styles.HelpDesc.Render("Scroll")),
	}

	helpContent := strings.Join(helpItems, " • ")
	
	// Apply styling and fit to width
	helpBar := m.styles.HelpBar.Width(m.width - 2).Render(helpContent)
	
	return helpBar
}
