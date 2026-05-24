package game

// Rotation represents a piece rotation state (0-3).
type Rotation int

const (
	R0 Rotation = iota
	R1
	R2
	R3
)

// Offset represents an (x, y) shift for wall kick attempts.
type Offset struct{ X, Y int }

// WallKickData returns the 5 wall kick offset attempts for a given
// piece type and rotation transition.
func WallKickData(pieceType int, from, to Rotation) []Offset {
	if pieceType == 0 { // I-piece
		return iOffsets[from][to]
	}
	return jlstzOffsets[from][to]
}

// JLSTZ wall kick offsets from SRS specification
var jlstzOffsets = map[Rotation]map[Rotation][]Offset{
	R0: {
		R1: {{X: 0, Y: 0}, {X: -1, Y: 0}, {X: -1, Y: -1}, {X: 0, Y: 2}, {X: -1, Y: 2}},
	},
	R1: {
		R0: {{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: -2}, {X: 1, Y: -2}},
		R2: {{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: -2}, {X: 1, Y: -2}},
	},
	R2: {
		R1: {{X: 0, Y: 0}, {X: -1, Y: 0}, {X: -1, Y: -1}, {X: 0, Y: 2}, {X: -1, Y: 2}},
		R3: {{X: 0, Y: 0}, {X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: 2}, {X: 1, Y: 2}},
	},
	R3: {
		R2: {{X: 0, Y: 0}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: -2}, {X: -1, Y: -2}},
		R0: {{X: 0, Y: 0}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: -2}, {X: -1, Y: -2}},
	},
}

// I-piece wall kick offsets from SRS specification
var iOffsets = map[Rotation]map[Rotation][]Offset{
	R0: {
		R1: {{X: 0, Y: 0}, {X: -2, Y: 0}, {X: 1, Y: 0}, {X: -2, Y: 1}, {X: 1, Y: -2}},
	},
	R1: {
		R0: {{X: 0, Y: 0}, {X: 2, Y: 0}, {X: -1, Y: 0}, {X: 2, Y: -1}, {X: -1, Y: 2}},
		R2: {{X: 0, Y: 0}, {X: -1, Y: 0}, {X: 2, Y: 0}, {X: -1, Y: -2}, {X: 2, Y: 1}},
	},
	R2: {
		R1: {{X: 0, Y: 0}, {X: 1, Y: 0}, {X: -2, Y: 0}, {X: 1, Y: -1}, {X: -2, Y: 2}},
		R3: {{X: 0, Y: 0}, {X: 2, Y: 0}, {X: -1, Y: 0}, {X: 2, Y: 1}, {X: -1, Y: -2}},
	},
	R3: {
		R2: {{X: 0, Y: 0}, {X: -2, Y: 0}, {X: 1, Y: 0}, {X: -2, Y: -1}, {X: 1, Y: 2}},
		R0: {{X: 0, Y: 0}, {X: 1, Y: 0}, {X: -2, Y: 0}, {X: 1, Y: 2}, {X: -2, Y: -1}},
	},
}
