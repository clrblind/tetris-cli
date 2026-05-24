package game

import "time"

const (
	LinesPerLevel = 10
)

var scoreMultipliers = map[int]int{
	1: 40,
	2: 100,
	3: 300,
	4: 1200,
}

// ScoreManager manages the player's score and level.
type ScoreManager struct {
	Score        int
	LinesCleared int
	Level        int
}

// NewScoreManager creates a new score manager.
func NewScoreManager() *ScoreManager {
	return &ScoreManager{
		Level: 1,
	}
}

// AddLines adds the score for the number of lines cleared.
func (s *ScoreManager) AddLines(lines int) {
	if lines <= 0 {
		return
	}
	s.LinesCleared += lines
	s.Score += scoreMultipliers[lines] * s.Level
	s.Level = 1 + (s.LinesCleared / LinesPerLevel)
}

// FallSpeed returns the duration for a piece to fall one step.
func (s *ScoreManager) FallSpeed() time.Duration {
	ms := 1000 - (s.Level-1)*50
	if ms < 100 {
		ms = 100
	}
	return time.Duration(ms) * time.Millisecond
}
