package game

import (
	"math/rand"
	"time"
)

// BagRandomizer implements the 7-bag randomizer for Tetris pieces.
// It shuffles all 7 piece types and deals from the bag, refilling when empty.
type BagRandomizer struct {
	bag []int
	rng *rand.Rand
}

// NewBagRandomizer creates a new BagRandomizer with a seeded RNG.
func NewBagRandomizer() *BagRandomizer {
	return &BagRandomizer{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Next returns the next piece type from the bag.
func (br *BagRandomizer) Next() int {
	if len(br.bag) == 0 {
		br.bag = []int{0, 1, 2, 3, 4, 5, 6}
		br.rng.Shuffle(len(br.bag), func(i, j int) {
			br.bag[i], br.bag[j] = br.bag[j], br.bag[i]
		})
	}
	typ := br.bag[len(br.bag)-1]
	br.bag = br.bag[:len(br.bag)-1]
	return typ
}

// defaultRandomizer is the package-level bag randomizer used by NewPiece.
var defaultRandomizer = NewBagRandomizer()
