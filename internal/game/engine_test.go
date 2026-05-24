package game

import (
	"testing"
)

func TestEngineTick(t *testing.T) {
	e := NewEngine()
	initialY := e.CurrentPiece.Pos.Y
	
	locked := e.Tick()
	if locked {
		t.Error("Piece should not be locked after first tick")
	}
	if e.CurrentPiece.Pos.Y != initialY+1 {
		t.Errorf("Expected piece to move down, got Y=%d", e.CurrentPiece.Pos.Y)
	}
}

func TestEngineMoveAndRotate(t *testing.T) {
	e := NewEngine()
	initialX := e.CurrentPiece.Pos.X
	
	e.MoveLeft()
	if e.CurrentPiece.Pos.X != initialX-1 {
		t.Errorf("Expected piece to move left, got X=%d", e.CurrentPiece.Pos.X)
	}
	
	e.MoveRight()
	if e.CurrentPiece.Pos.X != initialX {
		t.Errorf("Expected piece to move back right, got X=%d", e.CurrentPiece.Pos.X)
	}
}

func TestEngineGameOver(t *testing.T) {
	e := NewEngine()
	// Fill the top center of the board but leave a hole so lines aren't cleared
	for x := 0; x < BoardWidth-1; x++ {
		for y := 0; y < 5; y++ {
			e.Board[y][x] = 1
		}
	}
	
	// Spawning next piece should trigger game over if it's wide enough
	// Or just manually set a piece that will collide
	e.NextPiece = Piece{
		Shape: [][]int{{1, 1, 1, 1}},
		Pos:   Position{X: 0, Y: 0},
	}
	
	e.lockPiece()
	if !e.GameOver {
		t.Error("Expected game over when top area is occupied")
	}
}
