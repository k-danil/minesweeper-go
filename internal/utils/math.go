package utils

import "iter"

func Clamp(a, min, max int) int {
	if a < min {
		a = min
	} else if a >= max {
		a = max - 1
	}
	return a
}

// MatrixIterator returns a function that iterates over the matrix
// (-1, -1) (-1, 0) (-1, 1)
// ( 0, -1) ( 0, 0) ( 0, 1)
// ( 1, -1) ( 1, 0) ( 1, 1)
func MatrixIterator() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for i := -1; i <= 1; i++ {
			for j := -1; j <= 1; j++ {
				if !yield(i, j) {
					return
				}
			}
		}
	}
}
