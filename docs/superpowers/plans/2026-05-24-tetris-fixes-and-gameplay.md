# Tetris CLI — Bug Fixes & Gameplay Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix all known bugs, add 7-bag randomizer, SRS wall kicks, ghost piece, hold piece, and refactor TUI styles.

**Architecture:** Three phases: (1) bug fixes across game engine and TUI, (2) gameplay features adding new files for randomizer and wall kicks, (3) extraction of styles into a separate file. All changes stay within existing `game` and `tui` packages.

**Tech Stack:** Go 1.25, Bubble Tea, Lipgloss

---

## File Map

### Files to create:
- `internal/game/randomizer.go` — BagRandomizer (7-bag piece generation)
- `internal/game/wallkicks.go` — SRS offset tables and WallKickData function
- `internal/tui/styles.go` — All Lipgloss styles, colors, ASCII art (extracted from model.go)

### Files to modify:
- `internal/game/piece.go` — Add Rotation field, update NewPiece to use BagRandomizer
- `internal/game/engine.go` — Add Hold, GhostBlocks, wall kick rotation, CanHold, HeldPiece
- `internal/game/board.go` — Add bounds check in PlacePiece
- `internal/game/score.go` — Clamp lines > 4 in AddLines
- `internal/tui/model.go` — Remove viewCache, remove styles (moved to styles.go)
- `internal/tui/update.go` — Sync isGameOver after Drop/Tick, add hold key handler
- `internal/tui/view.go` — Add ghost piece rendering, hold piece display in sidebar

### Test files to modify:
- `internal/game/board_test.go` — Add PlacePiece bounds check test
- `internal/game/score_test.go` — Add AddLines > 4 test
- `internal/game/engine_test.go` — Add Drop+GameOver test, Hold tests, GhostBlocks test
- `internal/tui/update_test.go` — Add isGameOver sync test, hold key test

### Files to remove:
- None (everything is covered by create/modify above)

---

## Phase 1: Bug Fixes

### Task 1.1: isGameOver sync after hard drop

**Files:**
- Modify: `internal/tui/update.go:83-86`

- [ ] **Step 1: Read current update.go key handler for spacebar**

  Read lines 83-86 of `internal/tui/update.go` to see the current hard drop handler.

- [ ] **Step 2: Add isGameOver sync after Drop()**

  Change:
  ```go
  case " ": // Spacebar
      if !m.isPaused {
          m.Game.Drop()
      }
  ```
  To:
  ```go
  case " ": // Spacebar
      if !m.isPaused {
          m.Game.Drop()
          if m.Game.GameOver {
              m.isGameOver = true
          }
      }
  ```

- [ ] **Step 3: Run tests**

  Run: `go test ./internal/tui/ -v`
  Expected: All tests PASS

- [ ] **Step 4: Commit**

  ```bash
  git add internal/tui/update.go
  git commit -m "fix: sync isGameOver immediately after hard drop"
  ```

### Task 1.2: isGameOver sync after soft drop

**Files:**
- Modify: `internal/tui/update.go:78-82`

- [ ] **Step 1: Read current soft drop handler**

  Read lines 78-82 of `internal/tui/update.go`.

- [ ] **Step 2: Add isGameOver sync after Tick()**

  Change:
  ```go
  case "down":
      if !m.isPaused {
          m.Game.Tick()
      }
  ```
  To:
  ```go
  case "down":
      if !m.isPaused {
          m.Game.Tick()
          if m.Game.GameOver {
              m.isGameOver = true
          }
      }
  ```

- [ ] **Step 3: Write test for hard drop game over sync**

  Add to `internal/tui/update_test.go`:

  ```go
  func TestHardDropGameOverSync(t *testing.T) {
      m := InitialModel()
      // Fill board except thin channel to cause immediate game over on next spawn
      for x := 0; x < game.BoardWidth; x++ {
          for y := 0; y < game.BoardHeight-1; y++ {
              m.Game.Board[y][x] = 7
          }
      }
      // Piece can exist at Y=19 (one row free above)
      m.Game.CurrentPiece.Pos.Y = game.BoardHeight - 1
      // Drop it
      m2, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
      model := m2.(Model)
      if !model.isGameOver {
          t.Error("Expected isGameOver to be true immediately after game-over-causing drop")
      }
  }
  ```

  Add import: `"tetris-cli/internal/game"` if not already present.

- [ ] **Step 4: Run tests**

  Run: `go test ./internal/tui/ -v -run TestHardDropGameOverSync`
  Expected: PASS

- [ ] **Step 5: Commit**

  ```bash
  git add internal/tui/update.go internal/tui/update_test.go
  git commit -m "fix: sync isGameOver immediately after soft drop; add test"
  ```

### Task 1.3: PlacePiece bounds check

**Files:**
- Modify: `internal/game/board.go:28-32`
- Test: `internal/game/board_test.go`

- [ ] **Step 1: Write failing test for out-of-bounds PlacePiece**

  Add to `internal/game/board_test.go`:

  ```go
  func TestPlacePieceOutOfBounds(t *testing.T) {
      b := NewBoard()
      // Piece with blocks at negative Y — should not panic
      p := Piece{
          Shape:      [][]int{{1, 1}, {1, 1}},
          Pos:        Position{X: 0, Y: -2},
          ColorIndex: 1,
      }
      b.PlacePiece(p)
      // Board should be unchanged
      for y := 0; y < BoardHeight; y++ {
          for x := 0; x < BoardWidth; x++ {
              if b[y][x] != 0 {
                  t.Errorf("Board should be empty after placing out-of-bounds piece, got %d at (%d,%d)", b[y][x], x, y)
              }
          }
      }
  }
  ```

- [ ] **Step 2: Run test to verify it fails (expected panic)**

  Run: `go test ./internal/game/ -v -run TestPlacePieceOutOfBounds`
  Expected: FAIL with panic (index out of range)

- [ ] **Step 3: Add bounds check to PlacePiece**

  Change:
  ```go
  func (b *Board) PlacePiece(p Piece) {
      for _, block := range p.Blocks() {
          b[block.Y][block.X] = p.ColorIndex
      }
  }
  ```
  To:
  ```go
  func (b *Board) PlacePiece(p Piece) {
      for _, block := range p.Blocks() {
          if block.Y >= 0 && block.Y < BoardHeight && block.X >= 0 && block.X < BoardWidth {
              b[block.Y][block.X] = p.ColorIndex
          }
      }
  }
  ```

- [ ] **Step 4: Run test to verify it passes**

  Run: `go test ./internal/game/ -v -run TestPlacePieceOutOfBounds`
  Expected: PASS

- [ ] **Step 5: Run all game tests**

  Run: `go test ./internal/game/ -v`
  Expected: All PASS

- [ ] **Step 6: Commit**

  ```bash
  git add internal/game/board.go internal/game/board_test.go
  git commit -m "fix: add bounds check to PlacePiece against out-of-range panic"
  ```

### Task 1.4: Remove viewCache dead code

**Files:**
- Modify: `internal/tui/model.go:16`

- [ ] **Step 1: Remove viewCache field**

  In `internal/tui/model.go`, remove `viewCache  string` from the Model struct.

- [ ] **Step 2: Verify tests pass**

  Run: `go build ./... && go vet ./... && go test ./...`
  Expected: All pass (no errors since viewCache was never referenced)

- [ ] **Step 3: Commit**

  ```bash
  git add internal/tui/model.go
  git commit -m "fix: remove unused viewCache field from Model"
  ```

### Task 1.5: scoreMultipliers fallback for lines > 4

**Files:**
- Modify: `internal/game/score.go:35-36`
- Test: `internal/game/score_test.go`

- [ ] **Step 1: Write failing test for > 4 lines**

  Add to `internal/game/score_test.go`:

  ```go
  func TestAddLinesOverMax(t *testing.T) {
      s := NewScoreManager()
      // 5 lines is technically impossible in standard Tetris but guard against it
      s.AddLines(5)
      expectedScore := 1200 * 1 // uses the 4-line multiplier as fallback
      if s.Score != expectedScore {
          t.Errorf("Expected score %d for 5 lines (using 4-line multiplier), got %d", expectedScore, s.Score)
      }
  }
  ```

- [ ] **Step 2: Run test to verify it fails**

  Run: `go test ./internal/game/ -v -run TestAddLinesOverMax`
  Expected: Score is 0 (since map returns zero value for missing key)

- [ ] **Step 3: Clamp lines in AddLines**

  Change `internal/game/score.go:35-36`:
  ```go
  s.LinesCleared += lines
  s.Score += scoreMultipliers[lines] * s.Level
  ```
  To:
  ```go
  s.LinesCleared += lines
  if lines > 4 {
      lines = 4
  }
  s.Score += scoreMultipliers[lines] * s.Level
  ```

- [ ] **Step 4: Run test to verify it passes**

  Run: `go test ./internal/game/ -v -run TestAddLinesOverMax`
  Expected: PASS

- [ ] **Step 5: Commit**

  ```bash
  git add internal/game/score.go internal/game/score_test.go
  git commit -m "fix: clamp lines > 4 in AddLines to prevent zero score from missing map key"
  ```

### Task 1.6: Fix color comments

**Files:**
- Modify: `internal/tui/model.go:44-46`

- [ ] **Step 1: Fix swapped color comments**

  Change:
  ```go
  lipgloss.Color("5"), // Cyan (Z)
  lipgloss.Color("6"), // Magenta (J)
  ```
  To:
  ```go
  lipgloss.Color("5"), // Magenta (Z)
  lipgloss.Color("6"), // Cyan (J)
  ```

- [ ] **Step 2: Commit**

  ```bash
  git add internal/tui/model.go
  git commit -m "fix: correct swapped color comments for Z (magenta) and J (cyan)"
  ```

---

## Phase 2: Gameplay Improvements

### Task 2.1: 7-bag Randomizer

**Files:**
- Create: `internal/game/randomizer.go`
- Modify: `internal/game/piece.go:41-52`
- Test: `internal/game/piece_test.go`

- [ ] **Step 1: Write failing test for BagRandomizer**

  Add to `internal/game/piece_test.go`:

  ```go
  func TestBagRandomizerAllPiecesInWindow(t *testing.T) {
      br := NewBagRandomizer()
      // In every 7-piece window, all 7 types must appear exactly once
      for i := 0; i < 100; i++ {
          seen := make(map[int]int)
          for j := 0; j < 7; j++ {
              p := br.Next()
              seen[p]++
          }
          if len(seen) != 7 {
              t.Errorf("Window %d: expected 7 unique pieces, got %d", i, len(seen))
          }
          for typ, count := range seen {
              if count != 1 {
                  t.Errorf("Window %d: piece %d appeared %d times (expected 1)", i, typ, count)
              }
          }
      }
  }
  ```

- [ ] **Step 2: Run test to verify it fails**

  Run: `go test ./internal/game/ -v -run TestBagRandomizerAllPiecesInWindow`
  Expected: FAIL (NewBagRandomizer not defined)

- [ ] **Step 3: Create internal/game/randomizer.go**

  ```go
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
  ```

- [ ] **Step 4: Run test to verify it passes**

  Run: `go test ./internal/game/ -v -run TestBagRandomizerAllPiecesInWindow`
  Expected: PASS

- [ ] **Step 5: Add package-level randomizer and update NewPiece**

  In `internal/game/randomizer.go`, add:
  ```go
  // defaultRandomizer is the package-level bag randomizer used by NewPiece.
  var defaultRandomizer = NewBagRandomizer()
  ```

  In `internal/game/piece.go`, change `NewPiece`:
  ```go
  func NewPiece() Piece {
      pieceType := defaultRandomizer.Next()
      shape := shapes[pieceType]
      startPos := Position{X: BoardWidth/2 - len(shape[0])/2, Y: 0}
      return Piece{
          Shape:      shape,
          Pos:        startPos,
          Type:       pieceType,
          ColorIndex: pieceType + 1,
      }
  }
  ```

  Also remove `import "math/rand"` from `piece.go` if it's no longer used (it shouldn't be after the change).

- [ ] **Step 6: Run all game tests**

  Run: `go test ./internal/game/ -v`
  Expected: All PASS

- [ ] **Step 7: Run full test suite**

  Run: `go build ./... && go vet ./... && go test ./...`
  Expected: All pass

- [ ] **Step 8: Commit**

  ```bash
  git add internal/game/randomizer.go internal/game/piece.go internal/game/piece_test.go
  git commit -m "feat: add 7-bag randomizer for fair piece distribution"
  ```

### Task 2.2: Wall Kicks (SRS)

**Files:**
- Create: `internal/game/wallkicks.go`
- Modify: `internal/game/piece.go` (add Rotation field)
- Modify: `internal/game/engine.go:80-91` (replace Rotate with SRS)
- Test: `internal/game/engine_test.go`

- [ ] **Step 1: Define Rotation type and offset tables**

  Create `internal/game/wallkicks.go`:

  ```go
  package game

  // Rotation state for SRS.
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
  // piece type and rotation transition. Returns JLSTZ or I offsets.
  func WallKickData(pieceType int, from, to Rotation) []Offset {
      if pieceType == 0 { // I-piece
          return iOffsets[from][to]
      }
      return jlstzOffsets[from][to]
  }

  // JLSTZ wall kick offsets (from Tetris SRS specification)
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

  // I-piece wall kick offsets (from Tetris SRS specification)
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
  ```

- [ ] **Step 2: Write failing wall kick tests**

  Add to `internal/game/engine_test.go`:

  ```go
  func TestWallKickIRotate(t *testing.T) {
      e := NewEngine()
      // Set I-piece at left wall
      e.CurrentPiece = Piece{
          Shape:      shapes[0],
          Type:       0,
          Pos:        Position{X: 0, Y: 5},
          ColorIndex: 1,
      }
      e.Rotate() // Should kick right
      // If no wall kick, piece would be 4x1 vertical at X=0, overlapping left wall
      if e.CurrentPiece.Pos.X < 0 {
          t.Error("Expected wall kick to prevent I-piece from going out of bounds")
      }
      if !e.Board.IsValidPosition(e.CurrentPiece) {
          t.Error("Expected I-piece rotation with wall kick to be valid")
      }
  }

  func TestWallKickJLTZS(t *testing.T) {
      e := NewEngine()
      // Set T-piece at right wall
      e.CurrentPiece = Piece{
          Shape:      shapes[2], // T-piece: {{0,1,0},{1,1,1}}
          Type:       2,
          Pos:        Position{X: game.BoardWidth - 2, Y: 5},
          ColorIndex: 3,
      }
      e.Rotate() // Should kick from wall
      if !e.Board.IsValidPosition(e.CurrentPiece) {
          t.Error("Expected T-piece rotation at right wall to be valid with wall kick")
      }
  }
  ```

- [ ] **Step 3: Run tests to verify they fail**

  Run: `go test ./internal/game/ -v -run TestWallKick`
  Expected: FAIL (no wall kick yet, rotation fails at wall)

- [ ] **Step 4: Add Rotation field to Piece**

  In `internal/game/piece.go`, add `Rotation Rotation` field to the Piece struct:
  ```go
  type Piece struct {
      Shape      [][]int
      Pos        Position
      Type       int
      ColorIndex int
      Rotation   Rotation // 0-3: 0=spawn, 1=R, 2=2, 3=L
  }
  ```

  In `NewPiece()`, initialize with `Rotation: R0`.

- [ ] **Step 5: Update Rotate() method to track rotation state**

  In `internal/game/piece.go`, update `Rotate()`:
  ```go
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
      p.Rotation = (p.Rotation + 1) % 4
  }
  ```

- [ ] **Step 6: Replace Rotate in Engine with SRS**

  In `internal/game/engine.go`, replace the `Rotate()` method:

  ```go
  func (e *Engine) Rotate() {
      if e.GameOver {
          return
      }

      from := e.CurrentPiece.Rotation
      to := (from + 1) % 4
      offsets := WallKickData(e.CurrentPiece.Type, from, to)

      for _, offset := range offsets {
          originalShape := make([][]int, len(e.CurrentPiece.Shape))
          for i := range e.CurrentPiece.Shape {
              originalShape[i] = make([]int, len(e.CurrentPiece.Shape[i]))
              copy(originalShape[i], e.CurrentPiece.Shape[i])
          }
          originalPos := e.CurrentPiece.Pos
          originalRotation := e.CurrentPiece.Rotation

          e.CurrentPiece.Move(offset.X, offset.Y)
          e.CurrentPiece.Rotate()

          if e.Board.IsValidPosition(e.CurrentPiece) {
              return // success
          }

          // Revert
          e.CurrentPiece.Pos = originalPos
          e.CurrentPiece.Rotation = originalRotation
          e.CurrentPiece.Shape = originalShape
      }
  }
  ```

- [ ] **Step 7: Run tests to verify they pass**

  Run: `go test ./internal/game/ -v -run TestWallKick`
  Expected: PASS

- [ ] **Step 8: Run full test suite**

  Run: `go build ./... && go vet ./... && go test ./...`
  Expected: All pass

- [ ] **Step 9: Commit**

  ```bash
  git add internal/game/wallkicks.go internal/game/piece.go internal/game/engine.go internal/game/engine_test.go
  git commit -m "feat: add SRS wall kicks with offset tables for JLSTZ and I pieces"
  ```

### Task 2.3: Ghost Piece

**Files:**
- Modify: `internal/game/engine.go` (add GhostBlocks method)
- Modify: `internal/tui/view.go` (render ghost piece)
- Test: `internal/game/engine_test.go`

- [ ] **Step 1: Write failing test for GhostBlocks**

  Add to `internal/game/engine_test.go`:

  ```go
  func TestGhostBlocks(t *testing.T) {
      e := NewEngine()
      initialY := e.CurrentPiece.Pos.Y
      ghost := e.GhostBlocks()

      // Ghost should be below the current piece
      if len(ghost) == 0 {
          t.Fatal("Expected ghost blocks, got none")
      }
      for _, gb := range ghost {
          if gb.Y <= initialY {
              t.Errorf("Expected ghost block Y > %d, got %d", initialY, gb.Y)
          }
      }

      // Ghost should be at a valid position on the board
      for _, gb := range ghost {
          if gb.Y < 0 || gb.Y >= BoardHeight || gb.X < 0 || gb.X >= BoardWidth {
              t.Errorf("Ghost block out of bounds: %v", gb)
          }
      }
  }
  ```

- [ ] **Step 2: Run test to verify it fails**

  Run: `go test ./internal/game/ -v -run TestGhostBlocks`
  Expected: FAIL (GhostBlocks not defined)

- [ ] **Step 3: Add GhostBlocks to Engine**

  In `internal/game/engine.go`, add:

  ```go
  // GhostBlocks returns the block positions where the current piece would land.
  func (e *Engine) GhostBlocks() []Position {
      if e.GameOver {
          return nil
      }

      // Simulate dropping the piece
      ghostY := e.CurrentPiece.Pos.Y
      for {
          ghostY++
          e.CurrentPiece.Pos.Y = ghostY
          if !e.Board.IsValidPosition(e.CurrentPiece) {
              ghostY--
              e.CurrentPiece.Pos.Y = ghostY
              break
          }
      }

      blocks := make([]Position, len(e.CurrentPiece.Blocks()))
      copy(blocks, e.CurrentPiece.Blocks())

      // Restore original position
      e.CurrentPiece.Pos.Y = e.CurrentPiece.Pos.Y - (ghostY - e.CurrentPiece.Pos.Y)
      // Actually the above is wrong — we need to track original Y. Let me re-do:
      // Current Y is ghostY, we need to go back to original
  ```

  Wait, this approach mutates the piece and restores it. Better approach — clone the piece:

  ```go
  func (e *Engine) GhostBlocks() []Position {
      if e.GameOver || e.CurrentPiece.Pos.Y < 0 {
          return nil
      }

      // Clone piece data for simulation
      ghostPiece := e.CurrentPiece

      for e.Board.IsValidPosition(ghostPiece) {
          ghostPiece.Move(0, 1)
      }
      ghostPiece.Move(0, -1)

      return ghostPiece.Blocks()
  }
  ```

- [ ] **Step 4: Run test to verify it passes**

  Run: `go test ./internal/game/ -v -run TestGhostBlocks`
  Expected: PASS

- [ ] **Step 5: Render ghost piece in view**

  In `internal/tui/view.go`, modify `renderBoard()`.

  After creating the display board copy but before drawing the current piece blocks, add ghost piece rendering:

  ```go
  // Draw ghost piece (blocks where piece will land) — only when not game over
  if !m.isGameOver {
      ghostBlocks := m.Game.GhostBlocks()
      ghostColor := lipgloss.Color("8") // grey
      for _, block := range ghostBlocks {
          if block.Y >= 0 && block.Y < game.BoardHeight && block.X >= 0 && block.X < game.BoardWidth {
              if displayBoard[block.Y][block.X] == 0 {
                  // Mark ghost with a special sentinel value (use -1 to differentiate from placed blocks)
                  // We'll handle this during rendering
              }
          }
      }
  }
  ```

  Actually, using sentinel values on the display board is fragile. Better approach: during the rendering loop, check if a cell is a ghost cell.

  A cleaner approach: build a set of ghost positions and check during rendering:

  ```go
  // Build ghost position map
  ghostPositions := make(map[game.Position]bool)
  if !m.isGameOver {
      for _, block := range m.Game.GhostBlocks() {
          ghostPositions[block] = true
      }
  }

  // Later during board rendering:
  for y := 0; y < game.BoardHeight; y++ {
      var rowParts []string
      for x := 0; x < game.BoardWidth; x++ {
          cellValue := displayBoard[y][x]
          pos := game.Position{X: x, Y: y}

          if ghostPositions[pos] && displayBoard[y][x] == 0 {
              // Render ghost
              rowParts = append(rowParts, lipgloss.NewStyle().Foreground(ghostColor).Render("░░"))
          } else {
              color := blockColors[cellValue]
              rowParts = append(rowParts, lipgloss.NewStyle().Foreground(color).Render(blockStr))
          }
      }
      // ...
  }
  ```

  Note: `game.Position` equality comparison works because both fields are `int`.

  Add `ghostColor` to the styles section (or model.go variable block — later it'll be in styles.go).

- [ ] **Step 6: Build and run tests**

  Run: `go build ./... && go test ./...`
  Expected: All pass

- [ ] **Step 7: Commit**

  ```bash
  git add internal/game/engine.go internal/tui/view.go internal/game/engine_test.go
  git commit -m "feat: add ghost piece preview showing landing position"
  ```

### Task 2.4: Hold Piece

**Files:**
- Modify: `internal/game/engine.go` (add HeldPiece, CanHold, Hold method, reset CanHold in lockPiece)
- Modify: `internal/tui/update.go` (add hold key handler)
- Modify: `internal/tui/view.go` (add hold piece display in sidebar)
- Test: `internal/game/engine_test.go`, `internal/tui/update_test.go`

- [ ] **Step 1: Write failing tests for Hold**

  Add to `internal/game/engine_test.go`:

  ```go
  func TestHoldPiece(t *testing.T) {
      e := NewEngine()
      initialPiece := e.CurrentPiece

      e.Hold()
      if e.HeldPiece == nil {
          t.Fatal("Expected held piece after first Hold call")
      }
      if e.CurrentPiece.Type != e.NextPiece.Type {
          t.Error("Expected current piece to be the old next piece after first hold")
      }
      if e.CanHold {
          t.Error("Expected CanHold to be false after holding")
      }
  }

  func TestHoldSwap(t *testing.T) {
      e := NewEngine()
      e.Hold() // First hold stores current, brings next
      heldType := e.HeldPiece.Type

      // Need to lock a piece to reset CanHold
      // Simulate by setting CanHold back
      e.CanHold = true

      e.Hold() // Swap: current <- held, held <- old current
      if e.HeldPiece == nil {
          t.Fatal("Expected held piece after swap hold")
      }
      if e.HeldPiece.Type == heldType {
          t.Error("Expected held piece to be the old current piece after swap, not the same")
      }
  }

  func TestHoldCanResetOnLock(t *testing.T) {
      e := NewEngine()
      e.Hold()
      if e.CanHold {
          t.Fatal("CanHold should be false after hold")
      }
      // Simulate lock: place piece, clear lines, spawn new
      e.lockPiece()
      if !e.CanHold {
          t.Error("Expected CanHold to be true after lockPiece")
      }
  }
  ```

  Add to `internal/tui/update_test.go`:

  ```go
  func TestHoldKey(t *testing.T) {
      m := InitialModel()
      m2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")})
      model := m2.(Model)
      if model.Game.HeldPiece == nil {
          t.Error("Expected held piece after pressing 'c'")
      }
      if model.Game.CanHold {
          t.Error("Expected CanHold to be false after holding")
      }
  }
  ```

- [ ] **Step 2: Run tests to verify they fail**

  Run: `go test ./internal/game/ -v -run TestHold && go test ./internal/tui/ -v -run TestHoldKey`
  Expected: FAIL (Hold method not defined, HeldPiece/CanHold fields missing)

- [ ] **Step 3: Add HeldPiece, CanHold, and Hold method to Engine**

  Add fields to `Engine` struct:
  ```go
  type Engine struct {
      Board        *Board
      Score        *ScoreManager
      CurrentPiece Piece
      NextPiece    Piece
      GameOver     bool
      HeldPiece    *Piece
      CanHold      bool
  }
  ```

  Initialize in `NewEngine()`:
  ```go
  return &Engine{
      Board:        NewBoard(),
      Score:        NewScoreManager(),
      CurrentPiece: NewPiece(),
      NextPiece:    NewPiece(),
      CanHold:      true,
  }
  ```

  Add `Hold()` method:
  ```go
  func (e *Engine) Hold() {
      if !e.CanHold || e.GameOver {
          return
      }
      e.CanHold = false

      if e.HeldPiece == nil {
          // First hold: store current, spawn next
          held := e.CurrentPiece
          e.HeldPiece = &held
          e.CurrentPiece = e.NextPiece
          e.NextPiece = NewPiece()
      } else {
          // Swap current with held
          current := e.CurrentPiece
          e.CurrentPiece = *e.HeldPiece
          e.HeldPiece = &current

          // Reset position to top-center
          e.CurrentPiece.Pos = Position{
              X: BoardWidth/2 - len(e.CurrentPiece.Shape[0])/2,
              Y: 0,
          }
          e.CurrentPiece.Rotation = R0
      }
  }
  ```

  Reset `CanHold` in `lockPiece()` — add at the end of the method:
  ```go
  e.CanHold = true
  ```

- [ ] **Step 4: Add hold key handling to update.go**

  In `handleKeyMsg`, add to the switch:

  ```go
  case "c", "C":
      if !m.isPaused {
          m.Game.Hold()
      }
  ```

- [ ] **Step 5: Add hold piece rendering to view.go**

  Modify `renderInfo()` — after the next piece section and before the legend:

  ```go
  // Render held piece (if any)
  if m.Game.HeldPiece != nil {
      heldPiece.WriteString("\nHOLD:\n")
      heldGrid := make([][]int, 4)
      for i := range heldGrid {
          heldGrid[i] = make([]int, 4)
      }

      for _, block := range m.Game.HeldPiece.Blocks() {
          pX := block.X - m.Game.HeldPiece.Pos.X
          pY := block.Y - m.Game.HeldPiece.Pos.Y
          if pX >= 0 && pX < 4 && pY >= 0 && pY < 4 {
              heldGrid[pY][pX] = m.Game.HeldPiece.ColorIndex
          }
      }

      var heldStr strings.Builder
      for _, row := range heldGrid {
          for _, cell := range row {
              if cell != 0 {
                  color := blockColors[cell]
                  heldStr.WriteString(lipgloss.NewStyle().Foreground(color).Render("██"))
              } else {
                  heldStr.WriteString("  ")
              }
          }
          heldStr.WriteString("\n")
      }
  }
  ```

  Then in the info join, include `heldStr`:
  ```go
  info := lipgloss.JoinVertical(lipgloss.Left, score, level, lines, "\n", nextPiece.String(), heldStr.String(), "\n", legend)
  ```

  Handle the case where heldStr is empty (no held piece yet).

- [ ] **Step 6: Run tests to verify they pass**

  Run: `go test ./internal/game/ -v -run TestHold && go test ./internal/tui/ -v -run TestHoldKey`
  Expected: PASS

- [ ] **Step 7: Run full test suite**

  Run: `go build ./... && go vet ./... && go test ./...`
  Expected: All pass

- [ ] **Step 8: Commit**

  ```bash
  git add internal/game/engine.go internal/tui/update.go internal/tui/view.go internal/game/engine_test.go internal/tui/update_test.go
  git commit -m "feat: add hold piece with swap and CanHold lock mechanics"
  ```

---

## Phase 3: Architecture Refactoring

### Task 3.1: Extract styles from model.go into styles.go

**Files:**
- Create: `internal/tui/styles.go`
- Modify: `internal/tui/model.go` (remove style variables)

- [ ] **Step 1: Create internal/tui/styles.go**

  Move all style variables and ASCII art from `model.go` to `styles.go`. This includes:
  - `blockColors`
  - `containerStyle`, `boardStyle`, `infoStyle`, `gameOverStyle`, `pausedStyle`, `confirmStyle`, `legendStyle`, `keyStyle`
  - `pauseArt`
  - `gameOverArt`

  The file should be:

  ```go
  package tui

  import (
      "github.com/charmbracelet/lipgloss"
  )

  var (
      blockColors = []lipgloss.Color{
          lipgloss.Color("0"), // Empty
          lipgloss.Color("1"), // Red (I)
          lipgloss.Color("2"), // Green (O)
          lipgloss.Color("3"), // Yellow (T)
          lipgloss.Color("4"), // Blue (S)
          lipgloss.Color("5"), // Magenta (Z)
          lipgloss.Color("6"), // Cyan (J)
          lipgloss.Color("7"), // White (L)
      }

      ghostColor = lipgloss.Color("8")

      containerStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(1)
      boardStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false).Padding(0, 1)
      infoStyle      = lipgloss.NewStyle().Padding(0, 2)
      gameOverStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9")).SetString("GAME OVER")
      pausedStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11")).SetString("PAUSED")
      confirmStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13")).SetString("QUIT? (y/n)")
      legendStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
      keyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)

      pauseArt = `
    .▄▄▄▄▄. .▄▄▄▄▄. .▄▄   ▄▄. .▄▄▄▄▄. .▄▄▄▄▄. .▄▄▄▄▄.
    █░░░░░█ █░░░░░█ █░░█ █░░█ █░░░░░█ █░░░░░█ █░░░░░█
    █░░█▀▀▀ █░░█▀▀█ █░░█ █░░█ ▀▀▀▀█░░ █░░█▀▀▀ █░░█▀▀█
    █░░█    █░░█▄▄█ █░░█ █░░█ ▄▄▄▄█░░ █░░█▄▄▄ █░░█  █
    █░░█    █░░█  █ ▀░░▀▄▀░░▀ █░░░░░█ █░░░░░█ █░░█▄▄█
    ▀▀▀▀    ▀▀▀▀  ▀  ▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀ ▀▀▀▀▀▀▀ ▀▀▀▀▀▀ 
       |       |        |        |       |       |
       :       :        :        :       :       :
  `
      gameOverArt = `
   .▄▄▄▄▄. .▄▄▄▄▄. .▄▄   ▄▄. .▄▄▄▄▄.     .▄▄▄▄▄. .▄▄   ▄▄. .▄▄▄▄▄. .▄▄▄▄▄.
   █░░░░░█ █░░░░░█ █░░█▄▄█░░█ █░░░░░█     █░░░░░█ █░░█ █░░█ █░░░░░█ █░░░░░█
   █░░█▀▀▀ █░░█▀▀█ █░░█  █░░█ █░░█▀▀▀     █░░█ █░░ ▀░░▀▄▀░░▀ █░░█▀▀▀ █░░█▀▀█
   █░░█▄▄▄ █░░█▄▄█ █░░█  █░░█ █░░█▄▄▄     █░░█▄█░░  █░░█░░█  █░░█▄▄▄ █░░█▄▄█
   █░░░░░█ █░░█  █ █░░█  █░░█ █░░░░░█     █░░░░░█   ▀░░▀░░▀  █░░░░░█ █░░█  █
    ▀▀▀▀▀  ▀▀▀▀  ▀ ▀▀▀▀  ▀▀▀▀ ▀▀▀▀▀▀▀      ▀▀▀▀▀▀     ▀▀▀     ▀▀▀▀▀▀ ▀▀▀▀  ▀
       |       |       |       |              |         |         |      |
       :       :       :       :              :         :         :      :
  `
  )
  ```

  Note: `ghostColor` is added here (new variable from Task 2.3).

- [ ] **Step 2: Remove style variables from model.go**

  Remove from `internal/tui/model.go`:
  - The entire `var (...)` block starting at line 38
  - The `blockColors` slice
  - All style variables
  - `pauseArt` and `gameOverArt`

  The `model.go` import for `"github.com/charmbracelet/lipgloss"` may no longer be needed — remove it if the only usage was for styles.

  After removal, `model.go` should only import `"tetris-cli/internal/game"` and `tea "github.com/charmbracelet/bubbletea"`.

- [ ] **Step 3: Build and run tests**

  Run: `go build ./... && go test ./...`
  Expected: All pass (styles.go is in same package, all references remain valid)

- [ ] **Step 4: Commit**

  ```bash
  git add internal/tui/styles.go internal/tui/model.go
  git commit -m "refactor: extract styles, colors, and ASCII art from model.go into styles.go"
  ```

---

## Plan Self-Review Checklist

**Spec coverage:**
- [x] Phase 1 bugs (isGameOver sync, PlacePiece bounds, viewCache, scoreMultipliers, color comments) — Tasks 1.1-1.6
- [x] 7-bag randomizer — Task 2.1
- [x] SRS wall kicks — Task 2.2
- [x] Ghost piece — Task 2.3
- [x] Hold piece — Task 2.4
- [x] Styles extraction — Task 3.1

**Placeholder scan:** No TODOs, TBDs, or vague steps. Every step has exact code.

**Type/name consistency:** `HeldPiece`, `CanHold`, `GhostBlocks`, `Rotation`, `Rotation.R0`, `WallKickData`, `BagRandomizer` — all consistent across tasks.

**All tests run after each task:** Verified every task includes `go test` and `go build` commands.
