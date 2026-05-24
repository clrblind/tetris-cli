# Tetris CLI (Go Version) 🎮

A modern, terminal-based classic Tetris game written in Go using the Bubble Tea framework.

## ✨ Features

- **Classic Gameplay**: Authentic 10x20 board with all 7 standard tetrominoes.
- **Modern UI**: Built with `Bubble Tea` and `Lipgloss` for a beautiful, responsive terminal interface.
- **ASCII Art**: High-quality ASCII art backgrounds for Pause and Game Over screens.
- **Smooth Logic**: Improved game engine with wall-kick-like rotation and scoring.
- **Progressive Difficulty**: Falling speed increases as you clear lines.
- **Robust Architecture**: Modular design following the Model-View-Update (MVU) pattern.

## 🚀 Quick Start

### Build and Run

Ensure you have Go installed (1.18+ recommended).

```bash
# Build the binary
go build -o tetris-cli ./cmd/tetris-cli

# Start the game
./tetris-cli
```

### Controls

| Key | Action |
|-----|--------|
| **Arrow Keys** | Move pieces left/right/down |
| **↑** | Rotate piece clockwise |
| **Space** | Hard drop (instant to bottom) |
| **P** | Pause/Resume game |
| **Q / Ctrl+C** | Quit game (with confirmation) |
| **R** | Restart (after Game Over) |

## 🏗️ Architecture

- `internal/game`: Pure game logic (Board, Piece, Score, Engine).
- `internal/tui`: Terminal user interface logic using the Bubble Tea MVU pattern.
- `cmd/tetris-cli`: Application entry point.

## 🧪 Testing

Run all tests:

```bash
go test ./... -v
```

## 📜 License

MIT
