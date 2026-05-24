package tui

import (
	"tetris-cli/internal/game"
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the state of the TUI.
type Model struct {
	Game       *game.Engine
	isPaused   bool
	isGameOver bool
	isQuitting bool
	width      int
	height     int
}

// InitialModel creates the initial model for the TUI.
func InitialModel() Model {
	engine := game.NewEngine()
	return Model{
		Game:       engine,
		isPaused:   false,
		isGameOver: false,
	}
}

// Init is the first command that is run when the program starts.
func (m Model) Init() tea.Cmd {
	return doTick(m.Game.Score.FallSpeed())
}


