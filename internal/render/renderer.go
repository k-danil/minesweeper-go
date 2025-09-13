package render

import (
	"fmt"
	"minesweeper/internal/game"
)

const (
	cellCursor   = `[%s]`
	cellNoCursor = ` %s `
)

func RenderField(field *game.Field) {
	switch field.State {
	case game.Win:
		fmt.Println("You win!")
	case game.Loose:
		fmt.Println("You lose!")
	default:
		fmt.Print("                          \n")
	}

	for r := range field.Rows {
		for c := range field.Columns {
			t := field.GetTile(c, r)
			if field.Cursor.IsSelectedTile(c, r) {
				fmt.Printf(cellCursor, t.String())
			} else {
				fmt.Printf(cellNoCursor, t.String())
			}
		}
		fmt.Print("\n")
	}
	fmt.Printf("\x1b[%dA", field.Rows+1)
	fmt.Printf("\x1b[%dD", field.Columns)
}
