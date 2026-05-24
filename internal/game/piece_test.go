package game

import (
	"testing"
)

func TestNewPiece(t *testing.T) {
	p := NewPiece()
	if p.Pos.X < 0 || p.Pos.X >= BoardWidth {
		t.Errorf("Piece X position out of bounds: %d", p.Pos.X)
	}
	if p.Pos.Y != 0 {
		t.Errorf("Piece Y position should be 0, got %d", p.Pos.Y)
	}
	if len(p.Shape) == 0 {
		t.Error("Piece shape is empty")
	}
}

func TestBagRandomizerAllPiecesInWindow(t *testing.T) {
	br := NewBagRandomizer()
	// In every 7-piece window, all 7 types must appear exactly once
	for i := 0; i < 100; i++ {
		seen := make(map[int]int)
		for j := 0; j < 7; j++ {
			p := br.Next()
			seen[p]++
		}
		if len(seen) != 7 {
			t.Errorf("Window %d: expected 7 unique pieces, got %d", i, len(seen))
		}
		for typ, count := range seen {
			if count != 1 {
				t.Errorf("Window %d: piece %d appeared %d times (expected 1)", i, typ, count)
			}
		}
	}
}

func TestRotate(t *testing.T) {
	// Test rotation of I-piece (horizontal to vertical)
	p := Piece{
		Shape: [][]int{{1, 1, 1, 1}},
	}
	p.Rotate()
	if len(p.Shape) != 4 || len(p.Shape[0]) != 1 {
		t.Errorf("Expected rotated shape 4x1, got %dx%d", len(p.Shape), len(p.Shape[0]))
	}

	// Test rotation of O-piece (should stay 2x2)
	p = Piece{
		Shape: [][]int{{1, 1}, {1, 1}},
	}
	p.Rotate()
	if len(p.Shape) != 2 || len(p.Shape[0]) != 2 {
		t.Errorf("Expected O-piece to stay 2x2, got %dx%d", len(p.Shape), len(p.Shape[0]))
	}
}

func TestBlocks(t *testing.T) {
	p := Piece{
		Shape: [][]int{{1, 1}, {1, 1}},
		Pos:   Position{X: 5, Y: 5},
	}
	blocks := p.Blocks()
	if len(blocks) != 4 {
		t.Errorf("Expected 4 blocks, got %d", len(blocks))
	}
	expected := []Position{
		{5, 5}, {6, 5}, {5, 6}, {6, 6},
	}
	for i, pos := range blocks {
		if pos != expected[i] {
			t.Errorf("Block %d: expected %v, got %v", i, expected[i], pos)
		}
	}
}
