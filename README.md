# Minesweeper (terminal, Go)

A simple terminal Minesweeper implemented in Go. The game renders directly to your terminal, supports keyboard navigation (WASD or arrow keys), flagging, and a DFS-based "simple" opening mode.


## Requirements
- Go 1.24+ (module sets toolchain to go1.24.6)
- A Unix-like terminal (macOS/Linux). Windows users can run it under WSL or use a terminal that supports ANSI escape sequences and raw mode.


## Build
You can build with either the Makefile or the go tool.

- Using Makefile (outputs to ./build/minesweeper):
  ```sh
  make
  # or explicitly
  make minesweeper
  ```

- Using Go directly:
  ```sh
  go build -o build/minesweeper ./cmd/minesweeper.go
  ```


## Run
From the project root, after building:
```sh
./build/minesweeper [flags]
```
Or run without building (Go will build and run):
```sh
go run ./cmd/minesweeper.go [flags]
```

### Flags (startup options)
- `-rows int`     Row count (default: 15)
- `-columns int`  Column count (default: 15)
- `-percent int`  Percentage of tiles that are mines (default: 35). Values are clamped to [1, 100] internally.
- `-simple`       Use DFS to open tiles (default: true). When enabled, opening an empty tile recursively reveals its neighbors in a DFS manner.
- `-help`         Print usage and exit.

Examples:
```sh
# Start a 10x10 board with 20% mines
./build/minesweeper -rows 10 -columns 10 -percent 20

# Start a larger grid without DFS opening
./build/minesweeper -rows 20 -columns 30 -percent 25 -simple=false
```


## Controls
During the game, use:
- Move cursor: `W`/`A`/`S`/`D` or arrow keys
- Open tile: `Space`
- Place/remove flag: `F`
- Reset board: `R`
- Quit: `Q` or `Esc`

Status line shows "You win!" or "You lose!" when the game ends. Press `Space` to perform an action while playing; after win/lose, press `Space` to start a new game (the app resets on action in that state).


## Notes & Troubleshooting
- Terminal size: The renderer prints a status line and a `rows x columns` grid. Ensure your terminal window is wide/tall enough to display the board.
- Cursor visibility: The game hides the terminal cursor while running and restores it on exit.
- Raw mode: The input is read in raw mode. If the program crashes unexpectedly and your terminal looks odd (e.g., hidden cursor), try running `reset` (on macOS/Linux) to restore terminal state.
- Performance: Pure text rendering; should work well within typical terminal sizes.


## License
This project is distributed under the terms of the LICENSE file included in the repository.
