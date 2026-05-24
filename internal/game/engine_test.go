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

func TestGhostBlocks(t *testing.T) {
	e := NewEngine()
	initialY := e.CurrentPiece.Pos.Y
	ghost := e.GhostBlocks()

	if len(ghost) == 0 {
		t.Fatal("Expected ghost blocks, got none")
	}

	// Ghost should be below the current piece
	for _, gb := range ghost {
		if gb.Y <= initialY {
			t.Errorf("Expected ghost block Y > %d, got %d", initialY, gb.Y)
		}
	}

	// Ghost should be at valid positions on the board
	for _, gb := range ghost {
		if gb.Y < 0 || gb.Y >= BoardHeight || gb.X < 0 || gb.X >= BoardWidth {
			t.Errorf("Ghost block out of bounds: %v", gb)
		}
	}
}

func TestHoldFirstTime(t *testing.T) {
	e := NewEngine()
	initialType := e.CurrentPiece.Type
	nextType := e.NextPiece.Type

	e.Hold()
	if e.HeldPiece == nil {
		t.Fatal("Expected held piece after first Hold call")
	}
	if e.HeldPiece.Type != initialType {
		t.Errorf("Expected held piece type %d, got %d", initialType, e.HeldPiece.Type)
	}
	if e.CurrentPiece.Type != nextType {
		t.Errorf("Expected current piece to be the old next piece (type %d), got type %d", nextType, e.CurrentPiece.Type)
	}
	if e.CanHold {
		t.Error("Expected CanHold to be false after holding")
	}
}

func TestHoldSwap(t *testing.T) {
	e := NewEngine()
	e.Hold() // first hold
	firstHeldType := e.HeldPiece.Type
	currentAfterFirstHold := e.CurrentPiece.Type

	// Reset CanHold to simulate piece lock
	e.CanHold = true
	e.Hold() // swap: current <-> held

	if e.HeldPiece == nil {
		t.Fatal("Expected held piece after swap")
	}
	if e.CurrentPiece.Type != firstHeldType {
		t.Errorf("Expected current piece type %d (was held), got %d", firstHeldType, e.CurrentPiece.Type)
	}
	if e.HeldPiece.Type != currentAfterFirstHold {
		t.Errorf("Expected held piece type %d (was current), got %d", currentAfterFirstHold, e.HeldPiece.Type)
	}
}

func TestHoldDisallowsDoubleHold(t *testing.T) {
	e := NewEngine()
	e.Hold()
	currentAfterFirstHold := e.CurrentPiece.Type
	e.Hold() // should do nothing since CanHold is false
	if e.CurrentPiece.Type != currentAfterFirstHold {
		t.Error("Expected double hold to be ignored")
	}
}

func TestHoldCanResetOnLock(t *testing.T) {
	e := NewEngine()
	e.Hold()
	if e.CanHold {
		t.Fatal("CanHold should be false after hold")
	}

	for !e.Tick() {
	}

	if !e.CanHold {
		t.Error("Expected CanHold to be true after lockPiece")
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
