package main

import (
	"flag"
	"log"
	"minesweeper/internal/game"
	"minesweeper/internal/render"
	"os"
)

const (
	cursorUp    = "\x1b[A"
	cursorDown  = "\x1b[B"
	cursorLeft  = "\x1b[D"
	cursorRight = "\x1b[C"
	escape      = "\x1b"
)

func main() {
	rows := flag.Int("rows", 15, "Row count")
	columns := flag.Int("columns", 15, "Column count")
	percent := flag.Int("percent", 35, "Mines percent")
	simple := flag.Bool("simple", true, "Use DFS to open tiles")
	help := flag.Bool("help", false, "Shows usage")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	f := game.NewField(*columns, *rows, *percent, *simple)
	r, err := render.NewRenderer(*columns, *rows)
	if err != nil {
		log.Fatal("Error initializing renderer: ", err)
	}
	defer func() { _ = r.Close() }()

	cmd := make([]byte, 3)

	for {
		r.RenderField(f)

		var l int
		if l, err = os.Stdin.Read(cmd); err != nil {
			log.Fatal("Error reading from stdin: ", err)
		}

		switch string(cmd[:l]) {
		case "q", escape:
			return
		case "r":
			f.PushEvent(game.FieldReset)
		case "w", cursorUp:
			f.Cursor.Move(0, -1)
		case "s", cursorDown:
			f.Cursor.Move(0, 1)
		case "a", cursorLeft:
			f.Cursor.Move(-1, 0)
		case "d", cursorRight:
			f.Cursor.Move(1, 0)
		case "f":
			f.PushEvent(game.FieldFlag)
		case " ":
			f.PushEvent(game.FieldAction)
		}
	}
}
