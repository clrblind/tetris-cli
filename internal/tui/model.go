package tui

import (
	"tetris-cli/internal/game"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// Styles
var (
	// Colors for the blocks, 0 is empty
	blockColors = []lipgloss.Color{
		lipgloss.Color("0"), // Empty
		lipgloss.Color("1"), // Red (I)
		lipgloss.Color("2"), // Green (O)
		lipgloss.Color("3"), // Yellow (T)
		lipgloss.Color("4"), // Blue (S)
		lipgloss.Color("5"), // Cyan (Z)
		lipgloss.Color("6"), // Magenta (J)
		lipgloss.Color("7"), // White (L)
	}

	containerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)
	boardStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false).Padding(0, 1)
	infoStyle      = lipgloss.NewStyle().Padding(0, 2)
	gameOverStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).SetString("GAME OVER")
	pausedStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11")).SetString("PAUSED")
	confirmStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13")).SetString("QUIT? (y/n)")
	legendStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
	keyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)

	// ASCII Art (Scene Style)
	pauseArt = `
  .▄▄▄▄▄. .▄▄▄▄▄. .▄▄   ▄▄. .▄▄▄▄▄. .▄▄▄▄▄. .▄▄▄▄▄.
  █░░░░░█ █░░░░░█ █░░█ █░░█ █░░░░░█ █░░░░░█ █░░░░░█
  █░░█▀▀▀ █░░█▀▀█ █░░█ █░░█ ▀▀▀▀█░░ █░░█▀▀▀ █░░█▀▀█
  █░░█    █░░█▄▄█ █░░█ █░░█ ▄▄▄▄█░░ █░░█▄▄▄ █░░█  █
  █░░█    █░░█  █ ▀░░▀▄▀░░▀ █░░░░░█ █░░░░░█ █░░█▄▄█
  ▀▀▀▀    ▀▀▀▀  ▀  ▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀ ▀▀▀▀▀▀ 
     |       |        |        |       |       |
     :       :        :        :       :       :
`
	gameOverArt = `
 .▄▄▄▄▄. .▄▄▄▄▄. .▄▄   ▄▄. .▄▄▄▄▄.     .▄▄▄▄▄. .▄▄   ▄▄. .▄▄▄▄▄. .▄▄▄▄▄.
 █░░░░░█ █░░░░░█ █░░█▄▄█░░█ █░░░░░█     █░░░░░█ █░░█ █░░█ █░░░░░█ █░░░░░█
 █░░█▀▀▀ █░░█▀▀█ █░░█  █░░█ █░░█▀▀▀     █░░█ █░░ ▀░░▀▄▀░░▀ █░░█▀▀▀ █░░█▀▀█
 █░░█▄▄▄ █░░█▄▄█ █░░█  █░░█ █░░█▄▄▄     █░░█▄█░░  █░░█░░█  █░░█▄▄▄ █░░█▄▄█
 █░░░░░█ █░░█  █ █░░█  █░░█ █░░░░░█     █░░░░░█   ▀░░▀░░▀  █░░░░░█ █░░█  █
  ▀▀▀▀▀  ▀▀▀▀  ▀ ▀▀▀▀  ▀▀▀▀ ▀▀▀▀▀▀▀      ▀▀▀▀▀▀     ▀▀▀     ▀▀▀▀▀▀ ▀▀▀▀  ▀
     |       |       |       |              |         |         |      |
     :       :       :       :              :         :         :      :
`
)
