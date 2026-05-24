package game

const (
	BoardWidth  = 10
	BoardHeight = 20
)

// Position represents a coordinate on the board.
type Position struct {
	X, Y int
}

// Piece represents a Tetris piece (tetromino).
type Piece struct {
	Shape      [][]int
	Pos        Position
	Type       int
	ColorIndex int
}

var shapes = [][][]int{
	// I
	{{1, 1, 1, 1}},
	// O
	{{1, 1}, {1, 1}},
	// T
	{{0, 1, 0}, {1, 1, 1}},
	// S
	{{0, 1, 1}, {1, 1, 0}},
	// Z
	{{1, 1, 0}, {0, 1, 1}},
	// J
	{{1, 0, 0}, {1, 1, 1}},
	// L
	{{0, 0, 1}, {1, 1, 1}},
}

// NewPiece creates a new random Tetris piece.
func NewPiece() Piece {
	pieceType := defaultRandomizer.Next()
	shape := shapes[pieceType]
	// Position piece at top-center
	startPos := Position{X: BoardWidth/2 - len(shape[0])/2, Y: 0}
	return Piece{
		Shape:      shape,
		Pos:        startPos,
		Type:       pieceType,
		ColorIndex: pieceType + 1, // Use piece type as color index (1-7)
	}
}

// Rotate rotates the piece 90 degrees clockwise.
func (p *Piece) Rotate() {
	shape := p.Shape
	newShape := make([][]int, len(shape[0]))
	for i := range newShape {
		newShape[i] = make([]int, len(shape))
	}

	for i, row := range shape {
		for j, val := range row {
			newShape[j][len(shape)-1-i] = val
		}
	}
	p.Shape = newShape
}

// Move moves the piece by the given delta.
func (p *Piece) Move(dx, dy int) {
	p.Pos.X += dx
	p.Pos.Y += dy
}

// Blocks returns the absolute positions of the piece's blocks on the board.
func (p *Piece) Blocks() []Position {
	var blocks []Position
	for y, row := range p.Shape {
		for x, val := range row {
			if val != 0 {
				blocks = append(blocks, Position{X: p.Pos.X + x, Y: p.Pos.Y + y})
			}
		}
	}
	return blocks
}
