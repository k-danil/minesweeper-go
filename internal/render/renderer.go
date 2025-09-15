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
	escCursorPosition = "\x1b[%d;%dH"
	escClearScreen    = "\x1b[%dJ"
)

const (
	statusWin  = "You win!"
	statusLose = "You lose!"
)

type Renderer struct {
	term *terminal
	out  *bufio.Writer

	resetTerminal string
}

func NewRenderer(columns, rows int) (*Renderer, error) {
	term, err := initTerminal()
	if err != nil {
		err = fmt.Errorf("failed to initialize terminal: %w", err)
		return nil, err
	}
	fmt.Print(escHideCursor)

	resetTerminal := fmt.Sprintf(escCursorPosition, 0, 0)
	resetTerminal += fmt.Sprintf(escClearScreen, 0)
	fmt.Print(resetTerminal)

	return &Renderer{
		term:          term,
		out:           bufio.NewWriterSize(os.Stdout, columns*(rows+1)),
		resetTerminal: resetTerminal,
	}, nil
}

func (r *Renderer) RenderField(field *game.Field) {
	_, _ = r.out.WriteString(r.resetTerminal)

	switch field.State {
	case game.Win:
		_, _ = r.out.WriteString(statusWin)
	case game.Lose:
		_, _ = r.out.WriteString(statusLose)
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

	_ = r.out.Flush()
}

func (r *Renderer) Close() error {
	fmt.Print(escShowCursor)
	return r.term.restore()
}
