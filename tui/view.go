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
		m.styles.StatusBarKey.Render("Watching:"),
		m.styles.StatusBarValue.Render(m.watchDir))

	debounceItem := fmt.Sprintf("%s %s",
		m.styles.StatusBarKey.Render("Debounce:"),
		m.styles.StatusBarValue.Render(m.debounceDuration.String()))

	restartItem := fmt.Sprintf("%s %s",
		m.styles.StatusBarKey.Render("Restarts:"),
		m.styles.StatusBarValue.Render(fmt.Sprintf("%d", m.restartCount)))

	buildStatusItem := fmt.Sprintf("%s %s",
		m.styles.StatusBarKey.Render("Status:"),
		m.buildStatusString())

	// Join status items with separators
	separator := m.styles.LogEntryInfo.Render(" │ ")
	statusContent := strings.Join([]string{
		watchDirItem,
		debounceItem,
		restartItem,
		buildStatusItem,
	}, separator)

	// Apply styling and fit to width
	return m.styles.StatusBar.Width(m.width).Render(statusContent)
}

// renderLogPanel creates the scrollable log panel
func (m Model) renderLogPanel() string {
	var logContent string
	if len(m.logs) == 0 {
		logContent = m.styles.LogEntryInfo.Render("✨ Watching for file changes...")
	} else {
		visibleLogs := m.getVisibleLogs()
		var logLines []string
		for _, entry := range visibleLogs {
			logLines = append(logLines, m.formatLogEntry(entry))
		}
		logContent = strings.Join(logLines, "\n")
	}

	// Add scroll indicator if there are more logs
	totalLogs := len(m.logs)
	maxVisible := m.getMaxVisibleLogs()
	var scrollInfo string
	if totalLogs > maxVisible {
		start := m.logViewStart + 1
		end := m.logViewStart + len(m.getVisibleLogs())
		if end > totalLogs {
			end = totalLogs
		}
		scrollInfo = fmt.Sprintf(" %d-%d/%d ", start, end, totalLogs)
	}

	// Combine log content and scroll indicator
	logPanelHeight := m.height - 2 // Status bar and help bar
	return m.styles.LogPanel.
		Width(m.width - 2). // Account for border
		Height(logPanelHeight - 2).
		Render(lipgloss.JoinVertical(lipgloss.Left, logContent, lipgloss.PlaceHorizontal(m.width-2, lipgloss.Right, m.styles.LogEntryInfo.Render(scrollInfo))))
}

// renderHelpBar creates the bottom help bar
func (m Model) renderHelpBar() string {
	helpItems := []string{
		fmt.Sprintf("%s %s",
			m.styles.HelpKey.Render("q/ctrl+c"),
			m.styles.HelpDesc.Render("quit")),
		fmt.Sprintf("%s %s",
			m.styles.HelpKey.Render("r"),
			m.styles.HelpDesc.Render("restart")),
		fmt.Sprintf("%s %s",
			m.styles.HelpKey.Render("c"),
			m.styles.HelpDesc.Render("clear logs")),
		fmt.Sprintf("%s %s",
			m.styles.HelpKey.Render("↑/↓/j/k"),
			m.styles.HelpDesc.Render("scroll")),
		fmt.Sprintf("%s %s",
			m.styles.HelpKey.Render("PgUp/PgDn"),
			m.styles.HelpDesc.Render("fast scroll")),
	}

	separator := m.styles.LogEntryInfo.Render(" • ")
	helpContent := strings.Join(helpItems, separator)

	// Apply styling and fit to width
	return m.styles.HelpBar.Width(m.width).Render(helpContent)
}
