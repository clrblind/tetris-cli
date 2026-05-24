package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors for the blocks, 0 is empty
	ghostColor = lipgloss.Color("8")

	blockColors = []lipgloss.Color{
		lipgloss.Color("0"), // Empty
		lipgloss.Color("1"), // Red (I)
		lipgloss.Color("2"), // Green (O)
		lipgloss.Color("3"), // Yellow (T)
		lipgloss.Color("4"), // Blue (S)
		lipgloss.Color("5"), // Magenta (Z)
		lipgloss.Color("6"), // Cyan (J)
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
