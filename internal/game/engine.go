package game

// GameEngine manages the overall game state.
type Engine struct {
	Board        *Board
	Score        *ScoreManager
	CurrentPiece Piece
	NextPiece    Piece
	GameOver     bool
}

// NewEngine creates a new game engine.
func NewEngine() *Engine {
	return &Engine{
		Board:        NewBoard(),
		Score:        NewScoreManager(),
		CurrentPiece: NewPiece(),
		NextPiece:    NewPiece(),
	}
}

// Tick advances the game state by one step (piece falls).
// Returns true if the piece was locked.
func (e *Engine) Tick() bool {
	if e.GameOver {
		return false
	}

	// Try to move piece down
	e.CurrentPiece.Move(0, 1)

	if !e.Board.IsValidPosition(e.CurrentPiece) {
		// Can't move down, so move back and lock
		e.CurrentPiece.Move(0, -1)
		e.lockPiece()
		return true
	}
	return false
}

// lockPiece locks the current piece in place, clears lines, and spawns a new piece.
func (e *Engine) lockPiece() {
	e.Board.PlacePiece(e.CurrentPiece)
	lines := e.Board.ClearLines()
	e.Score.AddLines(lines)

	// Spawn next piece
	e.CurrentPiece = e.NextPiece
	e.NextPiece = NewPiece()

	// Check for game over
	if !e.Board.IsValidPosition(e.CurrentPiece) {
		e.GameOver = true
	}
}

// MoveLeft moves the current piece left if the position is valid.
func (e *Engine) MoveLeft() {
	if e.GameOver {
		return
	}
	e.CurrentPiece.Move(-1, 0)
	if !e.Board.IsValidPosition(e.CurrentPiece) {
		e.CurrentPiece.Move(1, 0) // move back
	}
}

// MoveRight moves the current piece right if the position is valid.
func (e *Engine) MoveRight() {
	if e.GameOver {
		return
	}
	e.CurrentPiece.Move(1, 0)
	if !e.Board.IsValidPosition(e.CurrentPiece) {
		e.CurrentPiece.Move(-1, 0) // move back
	}
}

// Rotate rotates the current piece if the new orientation is valid.
func (e *Engine) Rotate() {
	if e.GameOver {
		return
	}
	e.CurrentPiece.Rotate()
	if !e.Board.IsValidPosition(e.CurrentPiece) {
		// Rotate back
		for i := 0; i < 3; i++ {
			e.CurrentPiece.Rotate()
		}
	}
}

// Drop moves the piece down until it locks.
func (e *Engine) Drop() {
	if e.GameOver {
		return
	}
	for e.Board.IsValidPosition(e.CurrentPiece) {
		e.CurrentPiece.Move(0, 1)
	}
	e.CurrentPiece.Move(0, -1)
	e.lockPiece()
}
