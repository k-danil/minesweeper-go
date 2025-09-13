package game

import "strconv"

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
			str = strconv.Itoa(int(t.Adjacent))
		}
	}
	return str
}

func (t *Tile) PushEvent(event TileEvent) FieldEvent {
	switch {
	case t.State == Flagged && event == TileFlag:
		t.State = Closed

	case t.State == Closed && event == TileFlag:
		t.State = Flagged

	case t.State == Closed && event == TileOpen:
		t.State = Opened
		if t.Mine {
			return fieldTileMine
		} else {
			return fieldTileClean
		}
	}
	return fieldTileNoop
}
