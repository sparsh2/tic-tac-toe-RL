package board

import (
	"testing"
	"tictactoe/models"

	"github.com/stretchr/testify/assert"
)

func TestBoard(t *testing.T) {
	board := &Board{}
	board.Reset()
	board.Print()

	board.Apply(models.Move{X: 0, Y: 1, XO: models.X})
	board.Print()

	a := board.GetResult()
	assert.Nil(t, a)

	assert.False(t, board.Check(models.Move{X: 0, Y: 1, XO: models.O}))
	assert.False(t, board.Check(models.Move{X: 0, Y: 1, XO: models.X}))
	assert.False(t, board.Check(models.Move{X: 0, Y: 0, XO: models.None}))
	assert.False(t, board.Check(models.Move{X: 0, Y: 1, XO: models.None}))

	board.Apply(models.Move{X: 0, Y: 1, XO: models.O})
	assert.Equal(t, models.X, board.board[0][1])
}

func TestBoard_IsDecisive(t *testing.T) {
	b := &Board{}
	b.Reset()

	b.Apply(models.Move{X: 0, Y: 0, XO: models.X})
	b.Apply(models.Move{X: 1, Y: 0, XO: models.O})
	b.Apply(models.Move{X: 0, Y: 1, XO: models.X})
	b.Apply(models.Move{X: 1, Y: 1, XO: models.O})
	b.Apply(models.Move{X: 0, Y: 2, XO: models.X})

	assert.True(t, b.IsDecisive())
	assert.Equal(t, models.X, b.winner)

	b.Reset()

	b.board = [][]models.Turn{
		{models.None, models.None, models.O},
		{models.None, models.O, models.None},
		{models.O, models.None, models.None},
	}

	assert.True(t, b.IsDecisive())
	assert.Equal(t, models.O, b.winner)

	b.Reset()

	b.board = [][]models.Turn{
		{models.None, models.None, models.X},
		{models.None, models.O, models.None},
		{models.O, models.None, models.None},
	}

	assert.False(t, b.IsDecisive())
	assert.Equal(t, models.None, b.winner)

	b.Reset()

	b.board = [][]models.Turn{
		{models.None, models.None, models.X},
		{models.None, models.O, models.X},
		{models.O, models.None, models.X},
	}

	assert.True(t, b.IsDecisive())
	assert.Equal(t, models.X, b.winner)
}

func TestBoard_GetResult(t *testing.T) {
	b := &Board{}
	b.Reset()

	b.Apply(models.Move{X: 0, Y: 0, XO: models.X})
	b.Apply(models.Move{X: 1, Y: 0, XO: models.O})
	b.Apply(models.Move{X: 0, Y: 1, XO: models.X})
	b.Apply(models.Move{X: 1, Y: 1, XO: models.O})
	b.Apply(models.Move{X: 0, Y: 2, XO: models.X})

	assert.True(t, b.IsDecisive())
	assert.Equal(t, models.X, b.winner)

	res := b.GetResult()
	assert.Equal(t, "n", res.Draw)

	b.Reset()

	b.Apply(models.Move{X: 0, Y: 0, XO: models.X})
	b.Apply(models.Move{X: 1, Y: 0, XO: models.O})
	b.Apply(models.Move{X: 0, Y: 1, XO: models.X})
	b.Apply(models.Move{X: 1, Y: 1, XO: models.O})

	assert.False(t, b.IsDecisive())
	assert.Equal(t, models.None, b.winner)

	res = b.GetResult()
	assert.Nil(t, res)
}
