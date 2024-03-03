package board

import (
	"fmt"
	"math/rand"
	"tictactoe/models"
	"tictactoe/models/result"
)

type IBoard interface {
	Print()
	Reset()
	Apply(models.IMove)
	Check(models.IMove) bool
	IsDone() bool
	GetResult()
	IsDecisive() bool
	String() string
	GetAllMoves() []models.Move
	RevertLastMove()
}

type boardHistoryUnit struct {
	move  models.Move
	board [][]models.Turn
}

type Board struct {
	history []*boardHistoryUnit
	turn    models.Turn
	board   [][]models.Turn
	winner  models.Turn
}

func (b *Board) Print() {
	for i := 0; i < 3; i++ {
		line := ""
		for j := 0; j < 3; j++ {
			line += " "
			if b.board[i][j] != models.None {
				line += b.board[i][j].String()
			} else {
				line += " "
			}
			line += " "
			if j != 2 {
				line += "|"
			} else {
				line += "\n"
			}
		}
		fmt.Print(line)
		if i != 2 {
			fmt.Print("-----------\n")
		}
	}
}

func (b *Board) PrintHistory() {
	for _, his := range b.history {
		for i := 0; i < 3; i++ {
			line := ""
			for j := 0; j < 3; j++ {
				line += " "
				if b.board[i][j] != models.None {
					line += his.board[i][j].String()
				} else {
					line += " "
				}
				line += " "
				if j != 2 {
					line += "|"
				} else {
					line += "\n"
				}
			}
			fmt.Print(line)
			if i != 2 {
				fmt.Print("-----------\n")
			}
		}
		fmt.Println()
	}
}

func (b *Board) String() string {
	s := ""

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			s += b.board[i][j].String()
		}
	}

	return s
}

func (b *Board) GetAllMoves() []models.Move {
	res := []models.Move{}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b.board[i][j] == models.None {
				res = append(res, models.Move{X: i, Y: j, XO: b.turn})
			}
		}
	}

	rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})

	return res
}

func (b *Board) GetTurn() models.Turn {
	return b.turn
}

func (b *Board) Reset() {
	b.turn = models.X
	b.board = [][]models.Turn{
		{models.None, models.None, models.None},
		{models.None, models.None, models.None},
		{models.None, models.None, models.None},
	}
	b.winner = models.None
	b.history = []*boardHistoryUnit{}
}

func (b *Board) Check(move models.Move) bool {
	if move.XO != b.turn || b.board[move.X][move.Y] != models.None {
		return false
	}

	return true
}

func (b *Board) Apply(move models.Move) {
	if move.XO != b.turn || b.board[move.X][move.Y] != models.None {
		return
	}

	b.turn = b.turn.Flip()
	b.board[move.X][move.Y] = move.XO

	hu := &boardHistoryUnit{
		move: move,
		board: [][]models.Turn{
			{models.None, models.None, models.None},
			{models.None, models.None, models.None},
			{models.None, models.None, models.None},
		},
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			hu.board[i][j] = b.board[i][j]
		}
	}

	b.history = append(b.history, hu)
}

func (b *Board) GetResult() *result.Result {
	decisive := b.IsDecisive()
	done := b.IsDone()

	if decisive {
		winner := models.O
		loser := models.X
		if b.turn != models.X {
			winner = models.X
			loser = models.O
		}
		return &result.Result{
			Draw:   "n",
			Winner: winner,
			Loser:  loser,
		}
	} else if done {
		return &result.Result{
			Draw: "draw",
		}
	} else {
		return nil
	}
}

func (b *Board) IsDone() bool {
	empty := false
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			empty = empty || b.board[i][j] == models.None
		}
	}

	return !empty
}

func (b *Board) isDecisiveCheckColumn(c int) bool {
	a := b.board[0][c]
	done := true

	for i := 0; i < 3; i++ {
		done = done && b.board[i][c] == a
	}

	if done && a != models.None {
		b.winner = a
		return done
	}

	return false
}

func (b *Board) isDecisiveCheckRow(r int) bool {
	a := b.board[r][0]
	done := true

	for i := 0; i < 3; i++ {
		done = done && b.board[r][i] == a
	}

	if done && a != models.None {
		b.winner = a
		return done
	}

	return false
}

func (b *Board) isDecisiveCheckDiag(isP bool) bool {
	a := b.board[0][0]
	if !isP {
		a = b.board[0][2]
	}
	done := true
	for i := 0; i < 3; i++ {
		if isP {
			done = done && b.board[i][i] == a
		} else {
			done = done && b.board[i][2-i] == a
		}
	}

	if done && a != models.None {
		b.winner = a
		return done
	}

	return false
}

func (b *Board) IsDecisive() bool {
	done := false

	done = done || b.isDecisiveCheckColumn(0)
	done = done || b.isDecisiveCheckColumn(1)
	done = done || b.isDecisiveCheckColumn(2)

	done = done || b.isDecisiveCheckRow(0)
	done = done || b.isDecisiveCheckRow(1)
	done = done || b.isDecisiveCheckRow(2)

	done = done || b.isDecisiveCheckDiag(true)
	done = done || b.isDecisiveCheckDiag(false)

	return done
}

func (b *Board) RevertLastMove() {
	if len(b.history) == 0 {
		return
	}

	his := b.history[len(b.history)-1]
	b.board[his.move.X][his.move.Y] = models.None
	b.turn = b.turn.Flip()
	b.history = b.history[:len(b.history)-1]
	b.winner = models.None
}
