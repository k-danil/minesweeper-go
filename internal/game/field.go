package game

import (
	"log"
	"math/rand/v2"
)

type FieldState int

const (
	Init FieldState = iota
	Generate
	Playing
	Win
	Loose
)

type Field struct {
	State       FieldState
	Columns     int
	Rows        int
	tiles       []Tile
	Cursor      Cursor
	mines       int
	remainTiles int
}

func NewField(columns, rows int, percent int) *Field {
	if percent >= 99 {
		percent = 99
	}
	if percent < 1 {
		percent = 1
	}
	return &Field{
		Columns: columns,
		Rows:    rows,
		tiles:   make([]Tile, columns*rows),
		mines:   columns * rows / 100 * percent,
		Cursor:  Cursor{columns, rows, 0, 0},
	}
}

func (f *Field) GetState() FieldState {
	currentState := f.State
	switch f.State {
	case Init:
		f.reset()
	case Loose:
		f.openTiles(true)
	case Win:
		f.openTiles(false)
	default:
	}
	return currentState
}

func (f *Field) AdvanceState(x, y int, useDFS bool) {
	t := f.GetTile(x, y)
	if t == nil {
		return
	}

	switch f.State {
	case Playing:
		if t.Mine {
			f.State = Loose
		} else {
			if useDFS {
				f.dfs(x, y)
			} else if t.Open() {
				f.remainTiles--
			}

			if f.remainTiles == 0 {
				f.State = Win
			}
		}
	default:
	}
}

func (f *Field) reset() {
	for i := range f.tiles {
		f.tiles[i] = Tile{}
	}
	f.remainTiles = len(f.tiles) - f.mines
}

func (f *Field) getSliceIndex(x, y int) int {
	return x + y*f.Columns
}

func (f *Field) Randomize() {
	cursorIndex := f.getSliceIndex(f.Cursor.GetPosition())
	for range f.mines {
		for {
			j := rand.IntN(f.Columns * f.Rows)
			if j == cursorIndex {
				continue
			}
			if !f.tiles[j].Mine {
				f.tiles[j].Mine = true
				break
			}
		}
	}
	f.calculateAdjacent()
}

func (f *Field) openTiles(onlyMines bool) {
	for i := range f.tiles {
		if onlyMines && !f.tiles[i].Mine {
			continue
		}
		f.tiles[i].State = Open
	}
}

func (f *Field) calculateAdjacent() {
	for r := range f.Rows {
		for c := range f.Columns {
			t := f.GetTile(c, r)
			if t == nil {
				log.Fatal("Tile is nil")
			}
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if i == 0 && j == 0 {
						continue
					}
					adj := f.GetTile(c+i, r+j)
					if adj != nil && adj.Mine {
						t.Adjacent++
					}
				}
			}
		}
	}
}

func (f *Field) GetTile(x, y int) *Tile {
	if x < 0 || x >= f.Columns || y < 0 || y >= f.Rows {
		return nil
	}

	return &f.tiles[f.getSliceIndex(x, y)]
}

func (f *Field) dfs(x, y int) {
	t := f.GetTile(x, y)
	if t == nil || t.Mine || t.State != Closed || !t.Open() {
		return
	}
	f.remainTiles--

	f.dfs(x+1, y)
	f.dfs(x-1, y)
	f.dfs(x, y+1)
	f.dfs(x, y-1)
}

type Cursor struct {
	limitX, limitY int
	x, y           int
}

func (c *Cursor) Move(x, y int) {
	c.x += x
	if c.x < 0 {
		c.x = 0
	}
	if c.x >= c.limitX {
		c.x = c.limitX - 1
	}
	c.y += y
	if c.y < 0 {
		c.y = 0
	}
	if c.y >= c.limitY {
		c.y = c.limitY - 1
	}
}

func (c *Cursor) GetPosition() (int, int) {
	return c.x, c.y
}

func (c *Cursor) UnderCursor(x, y int) bool {
	return c.x == x && c.y == y
}
