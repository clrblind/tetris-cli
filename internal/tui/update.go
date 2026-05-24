package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

// tickMsg is a message sent on every game tick.
type tickMsg time.Time

// handleTick advances the game state.
func (m Model) handleTick() (tea.Model, tea.Cmd) {
	if m.isPaused || m.isGameOver || m.isQuitting {
		return m, doTick(m.Game.Score.FallSpeed())
	}

	m.Game.Tick()

	if m.Game.GameOver {
		m.isGameOver = true
	}

	return m, doTick(m.Game.Score.FallSpeed())
}

// handleKeyMsg handles user key presses.
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit with Ctrl+C
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}

	// Handle quitting confirmation
	if m.isQuitting {
		switch msg.String() {
		case "y", "Y":
			return m, tea.Quit
		case "n", "N", "esc", "q":
			m.isQuitting = false
			return m, nil
		}
		return m, nil
	}

	// Handle game over state
	if m.isGameOver {
		switch msg.String() {
		case "r", "R":
			newModel := InitialModel()
			newModel.width = m.width
			newModel.height = m.height
			return newModel, newModel.Init()
		case "q", "Q":
			return m, tea.Quit
		}
		return m, nil
	}

	// Normal game keys
	switch msg.String() {
	case "c", "C":
		if !m.isPaused {
			m.Game.Hold()
		}
	case "q", "Q":
		m.isQuitting = true
		return m, nil
	case "p", "P":
		m.isPaused = !m.isPaused
	case "left":
		if !m.isPaused {
			m.Game.MoveLeft()
		}
	case "right":
		if !m.isPaused {
			m.Game.MoveRight()
		}
	case "up":
		if !m.isPaused {
			m.Game.Rotate()
		}
	case "down":
		if !m.isPaused {
			// Soft drop
			m.Game.Tick()
			if m.Game.GameOver {
				m.isGameOver = true
			}
		}
	case " ": // Spacebar
		if !m.isPaused {
			m.Game.Drop()
			if m.Game.GameOver {
				m.isGameOver = true
			}
		}
	}
	return m, nil
}

// Update handles messages and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tickMsg:
		return m.handleTick()
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	}
	return m, nil
}

// doTick creates a command that waits for the given duration and then sends a tickMsg.
func doTick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
