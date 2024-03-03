package models

import "fmt"

type Turn string

func (t Turn) String() string {
	return string(t)
}

func (t Turn) Flip() Turn {
	if t == X {
		return O
	} else if t == O {
		return X
	} else {
		return t
	}
}

type IMove interface {
	Print()
}

type Move struct {
	X  int
	Y  int
	XO Turn
}

func (m *Move) Print() {
	fmt.Printf("Turn %s, move (%d,%d)\n", m.XO, m.X, m.Y)
}

const (
	X    Turn = "X"
	O    Turn = "O"
	None Turn = "-"
)
