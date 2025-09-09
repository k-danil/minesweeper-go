package game

import "fmt"

type State int

const (
	Closed State = iota
	Open
	Flagged
)

type Tile struct {
	State    State
	Mine     bool
	Adjacent int
}

func (t *Tile) String() string {
	var str string
	switch t.State {
	case Closed:
		str = "."
	case Flagged:
		str = "F"
	case Open:
		if t.Mine {
			str = "*"
		} else {
			str = fmt.Sprintf("%d", t.Adjacent)
		}
	}
	return str
}

func (t *Tile) Flag() {
	switch t.State {
	case Flagged:
		t.State = Closed
	case Closed:
		t.State = Flagged
	default:
	}

	return
}

func (t *Tile) Open() bool {
	switch t.State {
	case Closed:
		t.State = Open
		return true
	default:
	}
	return false
}
