package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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
