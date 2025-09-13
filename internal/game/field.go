package game

import (
	"math/rand/v2"
	"minesweeper/internal/utils"
)

type FieldState uint8

const (
	Init FieldState = iota
	Playing
	Win
	Loose
)

type FieldEvent uint8

const (
	FieldReset FieldEvent = iota
	FieldFlag
	FieldAction

	fieldTileNoop
	fieldTileClean
	fieldTileMine
)

type Field struct {
	State       FieldState
	useDFS      bool
	Columns     int
	Rows        int
	tiles       []Tile
	Cursor      Cursor
	mines       int
	remainTiles int
}

func NewField(columns, rows int, percent int, useDFS bool) *Field {
	percent = utils.Clamp(percent, 1, 100)
	tileCount := columns * rows
	return &Field{
		Columns: columns,
		Rows:    rows,
		useDFS:  useDFS,
		tiles:   make([]Tile, tileCount),
		mines:   tileCount / 100 * percent,
		Cursor:  Cursor{columns, rows, 0, 0},
	}
}

func (f *Field) PushEvent(event FieldEvent) {
	switch {
	case event == FieldReset:
		f.reset()
		f.State = Init

	case event == FieldFlag && f.State == Playing:
		f.pushTileEvent(TileFlag)

	case event == FieldAction && (f.State == Win || f.State == Loose):
		f.PushEvent(FieldReset)

	case event == FieldAction && f.State == Init:
		f.randomize()
		f.State = Playing
		fallthrough

	case event == FieldAction && f.State == Playing:
		f.pushTileEvent(TileOpen)

	case event == fieldTileClean && f.State == Playing:
		f.remainTiles--
		if f.remainTiles == 0 {
			f.openTiles(false)
			f.State = Win
		}

	case event == fieldTileMine && f.State == Playing:
		f.openTiles(true)
		f.State = Loose

	}
}

func (f *Field) pushTileEvent(event TileEvent) {
	x, y := f.Cursor.GetPosition()

	if f.useDFS && event == TileOpen {
		f.dfs(x, y, true)
	} else {
		t := f.GetTile(x, y)
		if t == nil {
			return
		}
		f.PushEvent(t.PushEvent(event))
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

func (f *Field) randomize() {
	cursorIndex := f.getSliceIndex(f.Cursor.GetPosition())
	for range f.mines {
		for {
			j := rand.IntN(len(f.tiles) - 1)
			if j == cursorIndex || f.tiles[j].Mine {
				continue
			}

			f.tiles[j].Mine = true
			break
		}
	}

	for r := range f.Rows {
		for c := range f.Columns {
			t := f.GetTile(c, r)
			if t == nil || !t.Mine {
				continue
			}
			for i, j := range utils.GetMatrixIterator() {
				if i == 0 && j == 0 {
					continue
				}
				if adj := f.GetTile(c+i, r+j); adj != nil {
					if !adj.Mine {
						adj.Adjacent++
					}
				}
			}
		}
	}
}

func (f *Field) openTiles(onlyMines bool) {
	for i := range f.tiles {
		if onlyMines && !f.tiles[i].Mine {
			continue
		}
		f.tiles[i].State = Opened
	}
}

func (f *Field) GetTile(x, y int) *Tile {
	if x < 0 || x >= f.Columns || y < 0 || y >= f.Rows {
		return nil
	}

	return &f.tiles[f.getSliceIndex(x, y)]
}

func (f *Field) dfs(x, y int, first bool) {
	t := f.GetTile(x, y)
	if f.State != Playing || t == nil || (!first && t.Mine) {
		return
	}

	fieldEvent := t.PushEvent(TileOpen)
	f.PushEvent(fieldEvent)
	if t.Adjacent > 0 || fieldEvent != fieldTileClean {
		return
	}

	for i, j := range utils.GetMatrixIterator() {
		if i == 0 && j == 0 {
			continue
		}
		f.dfs(x+i, y+j, false)
	}
}

type Cursor struct {
	limitX, limitY int
	x, y           int
}

func (c *Cursor) Move(x, y int) {
	c.x += x
	c.x = utils.Clamp(c.x, 0, c.limitX)
	c.y += y
	c.y = utils.Clamp(c.y, 0, c.limitY)
}

func (c *Cursor) GetPosition() (x int, y int) {
	return c.x, c.y
}

func (c *Cursor) IsSelectedTile(x, y int) bool {
	return c.x == x && c.y == y
}
