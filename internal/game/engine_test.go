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
		Shape:    [][]int{{1, 1, 1, 1}},
		Pos:      Position{X: 0, Y: 0},
		Rotation: R0,
	}
	
	e.lockPiece()
	if !e.GameOver {
		t.Error("Expected game over when top area is occupied")
	}
}

func TestWallKickIRotateAtLeftWall(t *testing.T) {
	e := NewEngine()
	e.CurrentPiece = Piece{
		Shape:      shapes[0],
		Type:       0,
		Pos:        Position{X: 0, Y: 5},
		ColorIndex: 1,
		Rotation:   R0,
	}
	e.Rotate()
	if e.CurrentPiece.Pos.X < 0 {
		t.Error("Expected wall kick to prevent I-piece from going out of bounds at left wall")
	}
	if !e.Board.IsValidPosition(e.CurrentPiece) {
		t.Error("Expected I-piece rotation with wall kick to be valid at left wall")
	}
}

func TestWallKickTRotateAtRightWall(t *testing.T) {
	e := NewEngine()
	e.CurrentPiece = Piece{
		Shape:      shapes[2], // T-piece: {{0,1,0},{1,1,1}}
		Type:       2,
		Pos:        Position{X: BoardWidth - 2, Y: 5},
		ColorIndex: 3,
		Rotation:   R0,
	}
	e.Rotate()
	if !e.Board.IsValidPosition(e.CurrentPiece) {
		t.Error("Expected T-piece rotation at right wall to be valid with wall kick")
	}
}

func TestWallKickAllFiveOffsetsFail(t *testing.T) {
	e := NewEngine()
	// I-piece pinned against left wall and blocks — all 5 offsets should fail
	e.CurrentPiece = Piece{
		Shape:      shapes[0],
		Type:       0,
		Pos:        Position{X: 0, Y: 0},
		ColorIndex: 1,
		Rotation:   R0,
	}
	// Fill board to block all offset attempts
	for x := 0; x < 3; x++ {
		e.Board[0][x] = 1
		e.Board[1][x] = 1
		e.Board[2][x] = 1
	}
	originalX := e.CurrentPiece.Pos.X
	originalRotation := e.CurrentPiece.Rotation
	e.Rotate()
	if e.CurrentPiece.Pos.X != originalX || e.CurrentPiece.Rotation != originalRotation {
		t.Error("Expected piece to remain unchanged when all wall kick offsets fail")
	}
}
