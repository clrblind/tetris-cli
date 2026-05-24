package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"tetris-cli/internal/game"
)

func TestUpdateKeyMsg(t *testing.T) {
	m := InitialModel()

	// Test Pause
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	newModel := m2.(Model)
	if !newModel.isPaused {
		t.Error("Expected game to be paused after pressing 'p'")
	}

	m3, _ := newModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("p")})
	newModel = m3.(Model)
	if newModel.isPaused {
		t.Error("Expected game to be unpaused after pressing 'p' again")
	}

	// Test Quit confirmation
	m4, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	newModel = m4.(Model)
	if !newModel.isQuitting {
		t.Error("Expected isQuitting to be true after pressing 'q'")
	}

	_, cmd := newModel.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if cmd == nil {
		t.Error("Expected tea.Quit command after confirming with 'y'")
	}
}

func TestUpdateMovement(t *testing.T) {
	m := InitialModel()
	initialX := m.Game.CurrentPiece.Pos.X

	// Test Move Left
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	newModel := m2.(Model)
	if newModel.Game.CurrentPiece.Pos.X != initialX-1 {
		t.Errorf("Expected X to be %d, got %d", initialX-1, newModel.Game.CurrentPiece.Pos.X)
	}

	// Test Move Right
	m3, _ := newModel.Update(tea.KeyMsg{Type: tea.KeyRight})
	newModel = m3.(Model)
	if newModel.Game.CurrentPiece.Pos.X != initialX {
		t.Errorf("Expected X to be back to %d, got %d", initialX, newModel.Game.CurrentPiece.Pos.X)
	}
}

func TestHardDropGameOverSync(t *testing.T) {
	m := InitialModel()
	// Fill the top 2 rows almost completely (leave column 9 empty so ClearLines
	// doesn't clear them). This blocks any new piece from spawning.
	for x := 0; x < game.BoardWidth-1; x++ {
		for y := 0; y < 2; y++ {
			m.Game.Board[y][x] = 7
		}
	}
	// Move piece to row 3 so it starts in a valid position (rows 2+ are empty)
	m.Game.CurrentPiece.Pos.Y = 3
	// Hard drop triggers lockPiece → new piece can't spawn → game over
	m2, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	model := m2.(Model)
	if !model.isGameOver {
		t.Errorf("Expected isGameOver to be true immediately after game-over-causing hard drop")
	}
}
