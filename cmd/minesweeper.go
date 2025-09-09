package main

import (
	"flag"
	"log"
	"minesweeper/internal/game"
	"minesweeper/internal/render"
	"os"
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

	f := game.NewField(*columns, *rows, *percent)

	terminal, err := render.InitTerminal()
	if err != nil {
		log.Fatal("Error initializing terminal: ", err)
	}
	defer func() { _ = terminal.Restore() }()

	cmd := make([]byte, 3)

	for {
		if f.GetState() == game.Init {
			render.RenderField(f)
			f.State = game.Generate
		}

		var l int
		if l, err = os.Stdin.Read(cmd); err != nil {
			log.Fatal("Error reading from stdin: ", err)
		}

		switch string(cmd[:l]) {
		case "q", "\x1b":
			os.Exit(0)
		case "w", "\x1b[A":
			f.Cursor.Move(0, -1)
		case "s", "\x1b[B":
			f.Cursor.Move(0, 1)
		case "d", "\x1b[C":
			f.Cursor.Move(1, 0)
		case "a", "\x1b[D":
			f.Cursor.Move(-1, 0)
		case "f":
			t := f.GetTile(f.Cursor.GetPosition())
			t.Flag()
		case " ":
			switch f.GetState() {
			case game.Loose, game.Win:
				f.State = game.Init
			case game.Generate:
				f.Randomize()
				f.State = game.Playing
				fallthrough
			case game.Playing:
				x, y := f.Cursor.GetPosition()
				f.AdvanceState(x, y, *simple)
			default:
			}
		default:
		}
		render.RenderField(f)
	}
}
