package game

import "fmt"

type TileState uint8

const (
	Closed TileState = iota
	Opened
	Flagged
)

type TileEvent uint8

const (
	TileFlag TileEvent = iota
	TileOpen
)

type Tile struct {
	State    TileState
	Mine     bool
	Adjacent uint8
}

func (t *Tile) String() string {
	var str string
	switch t.State {
	case Closed:
		str = "."
	case Flagged:
		str = "F"
	case Opened:
		if t.Mine {
			str = "*"
		} else {
			str = fmt.Sprintf("%d", t.Adjacent)
		}
	}
	return str
}

func (t *Tile) PushEvent(event TileEvent) FieldEvent {
	switch t.State {
	case Flagged:
		if event == TileFlag {
			t.State = Closed
		}
	case Closed:
		switch event {
		case TileFlag:
			t.State = Flagged
		case TileOpen:
			t.State = Opened
			if t.Mine {
				return fieldTileMine
			} else {
				return fieldTileClean
			}
		}
	default:
	}
	return fieldTileNoop
}
