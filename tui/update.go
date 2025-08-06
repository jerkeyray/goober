package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		waitForLogMessage(m.logChan),
		tea.EnterAltScreen,
	)
}

// Update handles all the IO and updates the model accordingly
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "r":
			// Force restart - this will be handled by the parent application
			return m, forceRestartCmd()

		case "c":
			m.ClearLogs()
			return m, nil

		case "up", "k":
			m.scrollUp()
			return m, nil

		case "down", "j":
			m.scrollDown()
			return m, nil

		case "pgup":
			for i := 0; i < 10; i++ {
				m.scrollUp()
			}
			return m, nil

		case "pgdown":
			for i := 0; i < 10; i++ {
				m.scrollDown()
			}
			return m, nil
		}

	case LogMessage:
		m.AddLog(msg.Message, msg.Type)
		return m, waitForLogMessage(m.logChan)

	case ForceRestartMsg:
		// This message is handled by the TUI integration
		return m, nil

	case BuildStatusMsg:
		m.SetBuildStatus(msg.Status)
		return m, nil

	case RestartCountMsg:
		m.IncrementRestartCount()
		return m, nil
	}

	return m, nil
}

// Custom message types
type ForceRestartMsg struct{}
type BuildStatusMsg struct {
	Status BuildStatus
}
type RestartCountMsg struct{}

// Commands
func forceRestartCmd() tea.Cmd {
	return func() tea.Msg {
		return ForceRestartMsg{}
	}
}

func SendBuildStatus(status BuildStatus) tea.Cmd {
	return func() tea.Msg {
		return BuildStatusMsg{Status: status}
	}
}

func SendRestartCount() tea.Cmd {
	return func() tea.Msg {
		return RestartCountMsg{}
	}
}

// waitForLogMessage waits for log messages from the channel
func waitForLogMessage(logChan chan LogMessage) tea.Cmd {
	return func() tea.Msg {
		return <-logChan
	}
}
