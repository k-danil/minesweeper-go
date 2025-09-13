package render

import (
	"bufio"
	"fmt"
	"minesweeper-go/internal/game"
	"os"
)

const (
	escHideCursor     = "\x1b[?25l"
	escShowCursor     = "\x1b[?25h"
	escMoveCursorUp   = "\x1b[%dA"
	escMoveCursorLeft = "\x1b[%dD"
)

const (
	statusWin  = "You win!"
	statusLose = "You lose!"
	statusInit = "          "
)

type Renderer struct {
	term *terminal
	out  *bufio.Writer

	resetCursor string
}

func NewRenderer(columns, rows int) (*Renderer, error) {
	term, err := initTerminal()
	if err != nil {
		err = fmt.Errorf("failed to initialize terminal: %w", err)
		return nil, err
	}
	fmt.Print(escHideCursor)

	return &Renderer{
		term:        term,
		out:         bufio.NewWriterSize(os.Stdout, columns*(rows+1)),
		resetCursor: fmt.Sprintf(escMoveCursorUp+escMoveCursorLeft, rows+1, columns),
	}, nil
}

func (r *Renderer) RenderField(field *game.Field) {
	switch field.State {
	case game.Win:
		_, _ = r.out.WriteString(statusWin)
	case game.Lose:
		_, _ = r.out.WriteString(statusLose)
	default:
		_, _ = r.out.WriteString(statusInit)
	}
	_, _ = r.out.WriteString("\n")

	for pos := range field.Iterator() {
		t := field.GetTile(pos)
		if field.Cursor.Position == pos {
			_, _ = r.out.WriteString("[")
			_, _ = r.out.WriteString(t.String())
			_, _ = r.out.WriteString("]")
		} else {
			_, _ = r.out.WriteString(" ")
			_, _ = r.out.WriteString(t.String())
			_, _ = r.out.WriteString(" ")
		}
		if pos.X == field.Size.X-1 {
			_, _ = r.out.WriteString("\n")
		}
	}

	_, _ = r.out.WriteString(r.resetCursor)
	_ = r.out.Flush()
}

func (r *Renderer) Close() error {
	fmt.Print(escShowCursor)
	return r.term.restore()
}
