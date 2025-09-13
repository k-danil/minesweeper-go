package utils

import "iter"

type Vec2 struct {
	X, Y int
}

func (v Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{v.X + v2.X, v.Y + v2.Y}
}

func (v Vec2) Less(v2 Vec2) bool {
	return v.X < v2.X || v.Y < v2.Y
}

func (v Vec2) GreaterOrEqual(v2 Vec2) bool {
	return v.X >= v2.X || v.Y >= v2.Y
}

// Clamp returns a value clamped to the range [min, max]
func Clamp(a, min, max int) int {
	if a < min {
		a = min
	} else if a >= max {
		a = max - 1
	}
	return a
}

// AroundIterator returns a function that iterates over the matrix
// (-1, -1) (-1, 0) (-1, 1)
// ( 0, -1) ( 0, 0) ( 0, 1)
// ( 1, -1) ( 1, 0) ( 1, 1)
func AroundIterator(pos Vec2) iter.Seq[Vec2] {
	return func(yield func(Vec2) bool) {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if !yield(pos.Add(Vec2{i, j})) {
					return
				}
			}
		}
	}
}
