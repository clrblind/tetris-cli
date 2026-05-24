# Tetris CLI — Bug Fixes & Gameplay Improvements Design

**Date:** 2026-05-24  
**Status:** Draft  
**Version:** 1.0

---

## Overview

This document describes the design for fixing all known bugs in the Tetris CLI project, adding gameplay improvements (7-bag randomizer, wall kicks, ghost piece, hold piece), and refactoring the architecture for better modularity.

The project is a terminal-based Tetris game written in Go using Bubble Tea and Lipgloss.

---

## Phase 1: Bug Fixes

### 1.1 Game Over sync delay after Drop

**File:** `internal/tui/update.go`

**Problem:** After `m.Game.Drop()` sets `GameOver = true`, `m.isGameOver` is not synchronized until the next tick. This causes a one-frame visual glitch where the board renders with the colliding piece drawn on top.

**Fix:** After `m.Game.Drop()` in the spacebar handler, immediately check and sync:

```go
case " ":
    if !m.isPaused {
        m.Game.Drop()
        if m.Game.GameOver {
            m.isGameOver = true
        }
    }
```

**Also apply to soft drop (down arrow):** `m.Game.Tick()` can also trigger game over via `lockPiece()`. After the Tick call, check and sync similarly:

```go
case "down":
    if !m.isPaused {
        m.Game.Tick()
        if m.Game.GameOver {
            m.isGameOver = true
        }
    }
```

### 1.2 PlacePiece missing bounds check

**File:** `internal/game/board.go`

**Problem:** `PlacePiece` writes to `b[block.Y][block.X]` without checking bounds. If a block has Y < 0, this panics.

**Fix:** Add bounds check:

```go
func (b *Board) PlacePiece(p Piece) {
    for _, block := range p.Blocks() {
        if block.Y >= 0 && block.Y < BoardHeight &&
            block.X >= 0 && block.X < BoardWidth {
            b[block.Y][block.X] = p.ColorIndex
        }
    }
}
```

### 1.3 viewCache dead code

**File:** `internal/tui/model.go`, line 16

**Problem:** `viewCache string` field is declared but never read or written.

**Fix:** Remove the field from the Model struct.

### 1.4 scoreMultipliers fallback for lines > 4

**File:** `internal/game/score.go`, line 36

**Problem:** `scoreMultipliers[lines]` returns 0 for lines > 4 (key not in map).

**Fix:** Handle the edge case:

```go
func (s *ScoreManager) AddLines(lines int) {
    if lines <= 0 {
        return
    }
    s.LinesCleared += lines
    if lines > 4 {
        lines = 4
    }
    s.Score += scoreMultipliers[lines] * s.Level
    s.Level = 1 + (s.LinesCleared / LinesPerLevel)
}
```

### 1.5 Swapped color comments

**File:** `internal/tui/model.go`, lines 44-46

**Problem:** Comments say `// Cyan (Z)` and `// Magenta (J)` but ANSI colors are swapped (5=Magenta, 6=Cyan).

**Fix:** Correct comments to match actual ANSI colors:

```go
lipgloss.Color("5"), // Magenta (Z)
lipgloss.Color("6"), // Cyan (J)
lipgloss.Color("7"), // White (L)
```

---

## Phase 2: Gameplay Improvements

### 2.1 7-bag Randomizer

**File:** `internal/game/randomizer.go` (new)

Reduce streak variance by shuffling all 7 piece types and dealing from the bag before reshuffling.

**Design:**

```go
type BagRandomizer struct {
    bag []int
    rng *rand.Rand
}

func NewBagRandomizer() *BagRandomizer { ... }
func (br *BagRandomizer) Next() int { ... }
```

- `Next()`: if bag is empty, fill with `{0,1,2,3,4,5,6}` and Fisher-Yates shuffle. Return and pop the last element.
- Seed `rng` with `time.Now().UnixNano()` for non-deterministic output.
- Store as package-level singleton or inject into `NewPiece()`.
- Tests: verify all 7 pieces appear in every 7-deal window.

**Integration:** `NewPiece()` calls `randomizer.Next()` instead of `rand.Intn()`.

### 2.2 Wall Kicks (SRS)

**File:** `internal/game/wallkicks.go` (new), modify `internal/game/engine.go` and `internal/game/piece.go`

Implement Super Rotation System wall kicks.

**Design:**

1. Add `Rotation` field (int 0-3: 0=spawn, 1=R(R), 2=2, 3=L(R')) to `Piece`.
2. Offset tables for JLSTZ and I pieces (from Tetris wiki SRS spec):
   - `JLSZT Offsets`: 5 offsets per each of the 4 rotation states (0→R, R→2, 2→L, L→0).
   - `I Offsets`: separate table for I-piece.
3. Modify `Engine.Rotate()`:
   - Save current piece state.
   - For each offset in the wall kick table for the current (from → to) transition:
     - Shift piece by offset.
     - Rotate piece.
     - If position is valid, keep it.
     - Otherwise, revert and try next offset.
   - If all 5 offsets fail, keep original state (no rotation).

**Edge cases:** Test each piece type at left wall, right wall, and floor; test I-piece separately (different offset table). Verify that at least one offset succeeds, and that all 5 offsets failing correctly rejects the rotation.

### 2.3 Ghost Piece

**File:** `internal/tui/view.go`

Show a translucent preview of where the piece will land.

**Design:**

1. Add `func (e *Engine) GhostBlocks() []Position`:
   - Clone current piece shape (or create a temporary piece at the same position).
   - Move it down until it collides (same logic as Drop but without locking).
   - Return the blocks at that final position.
2. In `renderBoard()`, draw ghost blocks with `░` character in uniform grey BEFORE drawing the current piece blocks.
3. Ghost color: always `lipgloss.Color("8")` (ANSI bright black = grey) regardless of piece type, for a consistent dimmed appearance distinct from the active piece.

### 2.4 Hold Piece

**Files:** `internal/game/engine.go`, `internal/tui/model.go`, `internal/tui/update.go`, `internal/tui/view.go`

Swap the current piece with a held piece.

**Design:**

1. **Engine changes:**
   - Add `HeldPiece *Piece` and `CanHold bool` fields.
   - `func (e *Engine) Hold()`:
     - If `!CanHold || GameOver`, return.
     - Set `CanHold = false` (prevents double-hold until next lock).
     - If `HeldPiece == nil`: store current piece, spawn next piece as current.
     - Else: swap current piece with held piece, reset position of the new current piece to top-center.
     - On `lockPiece()`: reset `CanHold = true`.
2. **Model changes:**
   - No new fields needed (Model holds `*Engine`, which now has `HeldPiece`).
3. **Update changes:**
   - Bind hold to `"c"` or `"C"` key (case-insensitive, consistent with other keys).
   - `case "c", "C": if !m.isPaused { m.Game.Hold(); ... }`.
4. **View changes:**
   - In `renderInfo()`, add "HOLD:" section below the controls legend.
   - Only show if `HeldPiece != nil`.
   - Render the held piece in a small grid (same approach as Next Piece).

---

## Phase 3: Architecture Refactoring

### 3.1 New file structure

```
internal/
├── game/
│   ├── board.go       Board, IsValidPosition, PlacePiece, ClearLines
│   ├── engine.go      Engine, Tick, MoveLeft/Right, Rotate, Drop, Hold
│   ├── piece.go       Position, Piece, Rotate, Move, Blocks
│   ├── randomizer.go  BagRandomizer
│   ├── score.go       ScoreManager, AddLines, FallSpeed
│   └── wallkicks.go   SRS offset tables, WallKickData
└── tui/
    ├── model.go       Model struct, InitialModel, Init
    ├── styles.go      blockColors, styles, ASCII art (moved from model.go)
    ├── update.go      Update, handleTick, handleKeyMsg, doTick
    └── view.go        View, renderBoard, renderInfo, ghost piece
```

### 3.2 Extractions

- **`styles.go`**: Move all `var (...)` block from `model.go` into this new file. Import `lipgloss` and `game` packages.
- **`randomizer.go`**: Extract bag logic from `NewPiece()`. `NewPiece()` now calls `bag.Next()` via a package-level `DefaultRandomizer`.
- **`wallkicks.go`**: Independent module, no extraction needed — purely new code.

No changes to `go.mod` or imports required — all files stay within the same packages.

---

## Testing Strategy

### Phase 1 Tests
- Game Over sync: new test in `tui/update_test.go` — hard drop when board is nearly full, verify `isGameOver` immediately.
- PlacePiece bounds: test in `board_test.go` — call PlacePiece with piece at Y=-1, verify no panic and no board mutation.
- scoreMultipliers > 4: test in `score_test.go` — AddLines(5), verify score is calculated correctly (use 4-line multiplier).

### Phase 2 Tests
- 7-bag: verify each 7-deal window contains all 7 piece types. Run 1000 deals, verify distribution.
- Wall kicks: test each piece type at each rotation transition at board edges, verify kick succeeds in expected cases.
- Ghost: visual-only, tested via existing engine Tick/Drop logic — verify `GhostBlocks()` returns positions below current piece.
- Hold: test swap logic, test CanHold reset after lock, test hold with nil initial piece.

### Phase 3 Tests
- No behavioral changes — existing tests must still pass after refactoring.
- Verify all imports are correct.

---

## Implementation Order

```
Phase 1 (Bugs):
  1.1 isGameOver sync       [tui/update.go]
  1.2 PlacePiece bounds     [game/board.go]
  1.3 viewCache removal     [tui/model.go]
  1.4 scoreMultipliers      [game/score.go]
  1.5 color comments        [tui/model.go]

Phase 1 → Code Review → Merge

Phase 2 (Gameplay):
  2.1 7-bag randomizer      [game/randomizer.go, game/piece.go]
  2.2 Wall kicks (SRS)      [game/wallkicks.go, game/piece.go, game/engine.go]
  2.3 Ghost piece           [game/engine.go, tui/view.go]
  2.4 Hold piece            [game/engine.go, tui/update.go, tui/view.go]

Phase 2 → Code Review → Merge

Phase 3 (Refactoring):
  3.1 Extract styles.go     [tui/styles.go, tui/model.go]
  3.2 Verify all tests pass

Phase 3 → Code Review → Done
```
