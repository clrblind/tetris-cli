package game

import (
	"testing"
)

func TestScoreManager(t *testing.T) {
	s := NewScoreManager()

	// Initial state
	if s.Score != 0 || s.Level != 1 {
		t.Errorf("Initial score/level mismatch: %d/%d", s.Score, s.Level)
	}

	// Single line clear
	s.AddLines(1)
	if s.Score != 40 {
		t.Errorf("Expected score 40 for 1 line at level 1, got %d", s.Score)
	}

	// Level up
	s.AddLines(9) // Total 10 lines
	if s.Level != 2 {
		t.Errorf("Expected level 2 after 10 lines, got %d", s.Level)
	}

	// Score at level 2
	s.AddLines(1) // 1 line * 40 * 2 (multiplier * level)
	expectedScore := 40 + (40 * 2)
	if s.Score != expectedScore {
		t.Errorf("Expected score %d, got %d", expectedScore, s.Score)
	}
}

func TestFallSpeed(t *testing.T) {
	s := NewScoreManager()
	s1 := s.FallSpeed()
	
	s.Level = 2
	s2 := s.FallSpeed()
	
	if s2 >= s1 {
		t.Errorf("Expected faster fall speed at level 2: %v >= %v", s2, s1)
	}
}
