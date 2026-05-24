package game

import (
	"testing"
)

func TestIsValidPosition(t *testing.T) {
	b := NewBoard()
	p := Piece{
		Shape:      [][]int{{1, 1}, {1, 1}},
		Pos:        Position{X: 0, Y: 0},
		ColorIndex: 1,
	}

	// Valid position
	if !b.IsValidPosition(p) {
		t.Error("Expected (0,0) to be a valid position")
	}

	// Out of bounds (left)
	p.Pos.X = -1
	if b.IsValidPosition(p) {
		t.Error("Expected (-1,0) to be invalid (out of bounds left)")
	}

	// Out of bounds (bottom)
	p.Pos.X = 0
	p.Pos.Y = BoardHeight - 1
	if b.IsValidPosition(p) {
		t.Error("Expected bottom position to be invalid (out of bounds bottom)")
	}

	// Collision
	p.Pos.Y = 0
	b[1][1] = 1 // Occupy a cell
	if b.IsValidPosition(p) {
		t.Error("Expected collision to be invalid")
	}
}

func TestPlacePieceOutOfBounds(t *testing.T) {
	b := NewBoard()
	p := Piece{
		Shape:      [][]int{{1, 1}, {1, 1}},
		Pos:        Position{X: 0, Y: -2},
		ColorIndex: 1,
	}
	b.PlacePiece(p)
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if b[y][x] != 0 {
				t.Errorf("Board should be empty after placing out-of-bounds piece, got %d at (%d,%d)", b[y][x], x, y)
			}
		}
	}
}

func TestPlacePieceAndClearLines(t *testing.T) {
	b := NewBoard()
	
	// Fill the bottom line
	for x := 0; x < BoardWidth; x++ {
		b[BoardHeight-1][x] = 1
	}

	cleared := b.ClearLines()
	if cleared != 1 {
		t.Errorf("Expected 1 line cleared, got %d", cleared)
	}

	// Top line should be empty after shift
	if b[0][0] != 0 {
		t.Error("Expected top cell to be empty after clearing lines")
	}
}
