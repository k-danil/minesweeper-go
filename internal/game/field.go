package game

import (
	"iter"
	"math/rand/v2"
	"minesweeper-go/internal/utils"
	"slices"
)

type FieldState uint8

const (
	Init FieldState = iota
	Playing
	Win
	Lose
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
	State  FieldState
	Size   utils.Vec2
	Cursor Cursor

	useDFS      bool
	tiles       []Tile
	mines       int
	remainTiles int
}

func NewField(columns, rows int, percent int, useDFS bool) *Field {
	percent = utils.Clamp(percent, 1, 100)
	tileCount := columns * rows
	mines := tileCount / 100 * percent
	size := utils.Vec2{X: columns, Y: rows}
	return &Field{
		Size:        size,
		useDFS:      useDFS,
		tiles:       make([]Tile, tileCount),
		mines:       mines,
		remainTiles: tileCount - mines,
		Cursor:      Cursor{border: size},
	}
}

func (f *Field) PushEvent(event FieldEvent) {
	switch {
	case event == FieldReset:
		f.reset()
		f.State = Init

	case event == FieldFlag && f.State == Playing:
		f.pushTileEvent(TileFlag)

	case event == FieldAction && (f.State == Win || f.State == Lose):
		f.PushEvent(FieldReset)

	case event == FieldAction && f.State == Init:
		f.randomize()
		f.State = Playing
		f.PushEvent(event)

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
		f.State = Lose

	}
}

func (f *Field) pushTileEvent(event TileEvent) {
	if f.useDFS && event == TileOpen {
		f.dfs(f.Cursor.Position, true)
	} else {
		if t := f.GetTile(f.Cursor.Position); t != nil {
			f.PushEvent(t.PushEvent(event))
		}
	}
}

func (f *Field) reset() {
	for i := range f.tiles {
		f.tiles[i] = Tile{}
	}
	f.remainTiles = len(f.tiles) - f.mines
}

func (f *Field) getSliceIndex(pos utils.Vec2) int {
	return pos.X + pos.Y*f.Size.X
}

func (f *Field) randomize() {
	cursorIndex := f.getSliceIndex(f.Cursor.Position)
	aroundIndexes := make([]int, 0, 9)
	for nbr := range utils.AroundIterator(f.Cursor.Position) {
		if f.GetTile(nbr) == nil {
			continue
		}
		aroundIndexes = append(aroundIndexes, f.getSliceIndex(nbr))
	}

	for m := range f.mines {
		for {
			j := rand.IntN(len(f.tiles) - 1)
			if f.tiles[j].Mine {
				continue
			}

			if len(f.tiles)-m > 9 {
				if slices.Contains(aroundIndexes, j) {
					continue
				}
			} else if j == cursorIndex {
				continue
			}

			f.tiles[j].Mine = true
			break
		}
	}

	for pos := range f.Iterator() {
		t := f.GetTile(pos)
		if t == nil || !t.Mine {
			continue
		}
		for nbr := range utils.AroundIterator(pos) {
			if nbr != pos {
				if nbrTile := f.GetTile(nbr); nbrTile != nil && !nbrTile.Mine {
					nbrTile.Adjacent++
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

func (f *Field) Iterator() iter.Seq[utils.Vec2] {
	return func(yield func(utils.Vec2) bool) {
		for y := range f.Size.Y {
			for x := range f.Size.X {
				if !yield(utils.Vec2{X: x, Y: y}) {
					return
				}
			}
		}
	}
}

func (f *Field) GetTile(pos utils.Vec2) *Tile {
	if pos.Less(utils.Vec2{}) || pos.GreaterOrEqual(f.Size) {
		return nil
	}

	return &f.tiles[f.getSliceIndex(pos)]
}

func (f *Field) dfs(pos utils.Vec2, first bool) {
	t := f.GetTile(pos)
	if f.State != Playing || t == nil || (!first && t.Mine) {
		return
	}

	fieldEvent := t.PushEvent(TileOpen)
	f.PushEvent(fieldEvent)
	if t.Adjacent > 0 || fieldEvent != fieldTileClean {
		return
	}

	for nbr := range utils.AroundIterator(pos) {
		if nbr != pos {
			f.dfs(nbr, false)
		}
	}
}

type Cursor struct {
	Position, border utils.Vec2
}

func (c *Cursor) Move(vec2 utils.Vec2) {
	c.Position = c.Position.Add(vec2)
	c.Position = utils.Vec2{
		X: utils.Clamp(c.Position.X, 0, c.border.X),
		Y: utils.Clamp(c.Position.Y, 0, c.border.Y),
	}
}
