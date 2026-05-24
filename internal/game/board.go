package game

// Board represents the game board.
type Board [BoardHeight][BoardWidth]int

// NewBoard creates a new empty game board.
func NewBoard() *Board {
	return &Board{}
}

// IsValidPosition checks if a piece is in a valid position on the board.
// It checks for out-of-bounds and collisions with existing blocks.
func (b *Board) IsValidPosition(p Piece) bool {
	for _, block := range p.Blocks() {
		// Check if out of bounds
		if block.X < 0 || block.X >= BoardWidth || block.Y < 0 || block.Y >= BoardHeight {
			return false
		}
		// Check for collision with existing blocks on the board
		if b[block.Y][block.X] != 0 {
			return false
		}
	}
	return true
}

// PlacePiece places a piece onto the board.
func (b *Board) PlacePiece(p Piece) {
	for _, block := range p.Blocks() {
		if block.Y >= 0 && block.Y < BoardHeight && block.X >= 0 && block.X < BoardWidth {
			b[block.Y][block.X] = p.ColorIndex
		}
	}
}

// ClearLines clears completed lines from the board and returns the number of lines cleared.
func (b *Board) ClearLines() int {
	linesCleared := 0
	for y := 0; y < BoardHeight; y++ {
		isLineFull := true
		for x := 0; x < BoardWidth; x++ {
			if b[y][x] == 0 {
				isLineFull = false
				break
			}
		}

		if isLineFull {
			linesCleared++
			// Move all lines above this one down
			for row := y; row > 0; row-- {
				for col := 0; col < BoardWidth; col++ {
					b[row][col] = b[row-1][col]
				}
			}
			// Clear the top line
			for col := 0; col < BoardWidth; col++ {
				b[0][col] = 0
			}
		}
	}
	return linesCleared
}
