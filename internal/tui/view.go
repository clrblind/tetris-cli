package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"tetris-cli/internal/game"
)

const (
	SidebarWidth = 15
	MinScale     = 1
)

// renderBoard renders the current state of the game board.
func (m Model) renderBoard() string {
	// Calculate scaling
	scaleX := (m.width - SidebarWidth - 4) / (game.BoardWidth * 2)
	scaleY := (m.height - 4) / game.BoardHeight
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}
	if scale < MinScale {
		scale = MinScale
	}

	blockStr := ""
	for i := 0; i < scale; i++ {
		blockStr += "██"
	}

	// Create a temporary board to draw the current piece on
	displayBoard := *m.Game.Board
	if !m.isGameOver {
		for _, block := range m.Game.CurrentPiece.Blocks() {
			if block.Y >= 0 && block.Y < game.BoardHeight && block.X >= 0 && block.X < game.BoardWidth {
				displayBoard[block.Y][block.X] = m.Game.CurrentPiece.ColorIndex
			}
		}
	}

	// Calculate ghost piece positions
	ghostPositions := make(map[game.Position]bool)
	if !m.isGameOver {
		for _, block := range m.Game.GhostBlocks() {
			if block.Y >= 0 && block.Y < game.BoardHeight &&
				block.X >= 0 && block.X < game.BoardWidth &&
				displayBoard[block.Y][block.X] == 0 {
				ghostPositions[block] = true
			}
		}
	}

	var rows []string
	for y := 0; y < game.BoardHeight; y++ {
		var rowParts []string
		for x := 0; x < game.BoardWidth; x++ {
			pos := game.Position{X: x, Y: y}
			if ghostPositions[pos] {
				rowParts = append(rowParts, lipgloss.NewStyle().Foreground(ghostColor).Render("░░"))
			} else {
				color := blockColors[displayBoard[y][x]]
				rowParts = append(rowParts, lipgloss.NewStyle().Foreground(color).Render(blockStr))
			}
		}
		singleRow := strings.Join(rowParts, "")
		// Repeat row for vertical scaling
		for i := 0; i < scale; i++ {
			rows = append(rows, singleRow)
		}
	}

	boardStr := strings.Join(rows, "\n")

	// Overlay text
	var overlay string
	var bgChar string = " "
	var bgColor lipgloss.Color
	var hasBgColor bool

	if m.isGameOver {
		overlay = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(gameOverArt) +
			"\n\n" + legendStyle.Render("[R]estart  [Q]uit")
		bgChar = "░"
		bgColor = lipgloss.Color("233")
		hasBgColor = true
	} else if m.isQuitting {
		overlay = lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true).Render(
			"   .▄▄▄▄▄. .▄▄   ▄▄. .▄▄▄▄▄. .▄▄▄▄▄.\n"+
				"   █░░░░░█ █░░█ █░░█ █░░░░░█ ▀▀█░░█▀\n"+
				"   █░░█▀▀▀ ▀░░▀▄▀░░▀ █░░█▀▀▀    █░░█\n"+
				"   █░░█▄▄▄  █░░█░░█  █░░█▄▄▄    █░░█\n"+
				"   █░░░░░█  ▀░░▀░░▀  █░░░░░█    █░░█\n"+
				"    ▀▀▀▀▀     ▀▀▀     ▀▀▀▀▀     ▀▀▀ \n\n",
		) + confirmStyle.String()
		bgChar = " "
	} else if m.isPaused {
		overlay = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Render(pauseArt) +
			"\n\n" + legendStyle.Render("Press [P] to Resume")
		bgChar = "▒"
		bgColor = lipgloss.Color("234") // Darker grey
		hasBgColor = true
	}

	if overlay != "" {
		opts := []lipgloss.WhitespaceOption{
			lipgloss.WithWhitespaceChars(bgChar),
		}
		if hasBgColor {
			opts = append(opts, lipgloss.WithWhitespaceForeground(bgColor))
		}

		overlayBox := lipgloss.Place(game.BoardWidth*2*scale, game.BoardHeight*scale,
			lipgloss.Center, lipgloss.Center, overlay,
			opts...,
		)
		return boardStyle.Render(overlayBox)
	}

	return boardStyle.Render(boardStr)
}

// renderInfo renders the score, level, and next piece.
func (m Model) renderInfo() string {
	score := fmt.Sprintf("SCORE: %d", m.Game.Score.Score)
	level := fmt.Sprintf("LEVEL: %d", m.Game.Score.Level)
	lines := fmt.Sprintf("LINES: %d", m.Game.Score.LinesCleared)

	// Render next piece
	var nextPiece strings.Builder
	nextPiece.WriteString("NEXT:\n")

	// Create a small grid for the next piece
	grid := make([][]int, 4)
	for i := range grid {
		grid[i] = make([]int, 4)
	}

	// Draw the next piece onto the grid
	for _, block := range m.Game.NextPiece.Blocks() {
		// Normalize position to fit in the small grid
		pX := block.X - m.Game.NextPiece.Pos.X
		pY := block.Y - m.Game.NextPiece.Pos.Y
		if pX >= 0 && pX < 4 && pY >= 0 && pY < 4 {
			grid[pY][pX] = m.Game.NextPiece.ColorIndex
		}
	}

	for _, row := range grid {
		for _, cell := range row {
			if cell != 0 {
				color := blockColors[cell]
				nextPiece.WriteString(lipgloss.NewStyle().Foreground(color).Render("██"))
			} else {
				nextPiece.WriteString("  ")
			}
		}
		nextPiece.WriteString("\n")
	}

	// Legend
	legend := lipgloss.JoinVertical(lipgloss.Left,
		"\nCONTROLS:",
		fmt.Sprintf("%s Move", keyStyle.Render("← →")),
		fmt.Sprintf("%s Rotate", keyStyle.Render(" ↑ ")),
		fmt.Sprintf("%s Soft Drop", keyStyle.Render(" ↓ ")),
		fmt.Sprintf("%s Hard Drop", keyStyle.Render("Spc")),
		fmt.Sprintf("%s Pause", keyStyle.Render(" P ")),
		fmt.Sprintf("%s Quit", keyStyle.Render(" Q ")),
	)

	// Render held piece (if any)
	var heldPieceStr string
	if m.Game.HeldPiece != nil {
		var heldBuilder strings.Builder
		heldBuilder.WriteString("\nHOLD:\n")

		heldGrid := make([][]int, 4)
		for i := range heldGrid {
			heldGrid[i] = make([]int, 4)
		}
		for _, block := range m.Game.HeldPiece.Blocks() {
			pX := block.X - m.Game.HeldPiece.Pos.X
			pY := block.Y - m.Game.HeldPiece.Pos.Y
			if pX >= 0 && pX < 4 && pY >= 0 && pY < 4 {
				heldGrid[pY][pX] = m.Game.HeldPiece.ColorIndex
			}
		}
		for _, row := range heldGrid {
			for _, cell := range row {
				if cell != 0 {
					color := blockColors[cell]
					heldBuilder.WriteString(lipgloss.NewStyle().Foreground(color).Render("██"))
				} else {
					heldBuilder.WriteString("  ")
				}
			}
			heldBuilder.WriteString("\n")
		}
		heldPieceStr = heldBuilder.String()
	}

	info := lipgloss.JoinVertical(lipgloss.Left, score, level, lines, "\n", nextPiece.String(), heldPieceStr, "\n", legend)
	return infoStyle.Render(info)
}

// View renders the entire TUI.
func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}
	board := m.renderBoard()
	info := m.renderInfo()

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		containerStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, board, info)))
}
