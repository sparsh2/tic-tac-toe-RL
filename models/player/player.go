package player

import (
	"log"
	"math/rand"
	"tictactoe/models"
	"tictactoe/models/board"
)

type ImmutableBoard [][]models.Turn

type IPlayer interface {
	RegisterXO()
	Move(ImmutableBoard) models.IMove
	Result(string)
	GetStats()
}

type Human struct {
	totalGames int
	wins       int
	loses      int
}

type AgentMode int

const (
	PlayMode AgentMode = iota
	LearnMode
)

type AgentParameters struct {
	Mode                   AgentMode
	ExploratoryProbability float64
	LearningRate           float64
}

type Agent struct {
	totalGames int
	wins       int
	loses      int
	turn       models.Turn
	valueFn    map[string]float64
	params     *AgentParameters
	logger     *log.Logger
}

func NewAgent(valueFn map[string]float64, params *AgentParameters, logger *log.Logger, turn models.Turn) *Agent {
	return &Agent{
		turn:       turn,
		valueFn:    valueFn,
		logger:     logger,
		totalGames: 0,
		wins:       0,
		loses:      0,
		params:     params,
	}
}

func (a *Agent) RegisterXO(xo models.Turn) {
	a.turn = xo
}

func (a *Agent) Move(b board.Board) models.Move {
	if a.params.Mode == PlayMode {
		if a.turn == models.X {
			moves := b.GetAllMoves()
			if len(moves) == 0 {
				a.logger.Println("Got 0 available legal moves on the board")
				return models.Move{X: 0, Y: 0, XO: a.turn}
			}

			moveToPlay := moves[0]
			b.Apply(moveToPlay)
			maxVal, ok := a.valueFn[b.String()]
			if !ok {
				a.logger.Printf("Couldn't find the position in the map!")
				return models.Move{X: 0, Y: 0, XO: a.turn}
			}
			b.RevertLastMove()

			for _, mv := range moves {
				b.Apply(mv)
				val, ok := a.valueFn[b.String()]
				if !ok {
					a.logger.Printf("Couldn't find the position in the map!")
					return models.Move{X: 0, Y: 0, XO: a.turn}
				}
				if val > maxVal {
					moveToPlay = mv
					maxVal = val
				}
				b.RevertLastMove()
			}

			return moveToPlay
		} else {
			moves := b.GetAllMoves()
			if len(moves) == 0 {
				a.logger.Println("Got 0 available legal moves on the board")
				return models.Move{X: 0, Y: 0, XO: a.turn}
			}

			moveToPlay := moves[0]
			b.Apply(moveToPlay)
			minVal, ok := a.valueFn[b.String()]
			if !ok {
				a.logger.Printf("Couldn't find the position in the map!")
				return models.Move{X: 0, Y: 0, XO: a.turn}
			}
			b.RevertLastMove()

			for _, mv := range moves {
				b.Apply(mv)
				val, ok := a.valueFn[b.String()]
				if !ok {
					a.logger.Printf("Couldn't find the position in the map!")
					return models.Move{X: 0, Y: 0, XO: a.turn}
				}
				if val < minVal {
					moveToPlay = mv
					minVal = val
				}
				b.RevertLastMove()
			}

			return moveToPlay
		}

	} else {
		exploratory := true

		if rand.Float64() < a.params.ExploratoryProbability {
			exploratory = true
			a.logger.Println("exploring possibilities")
		} else {
			exploratory = false
			a.logger.Println("playing optimally")
		}
		if a.turn == models.X {
			if exploratory {
				moves := b.GetAllMoves()
				ind := rand.Intn(len(moves))
				return moves[ind]
			} else {
				moves := b.GetAllMoves()
				if len(moves) == 0 {
					a.logger.Println("Got 0 available legal moves on the board")
					return models.Move{X: 0, Y: 0, XO: a.turn}
				}

				moveToPlay := moves[0]
				b.Apply(moveToPlay)
				maxVal, ok := a.valueFn[b.String()]
				if !ok {
					a.logger.Printf("Couldn't find the position in the map!")
					return models.Move{X: 0, Y: 0, XO: a.turn}
				}
				b.RevertLastMove()

				for _, mv := range moves {
					b.Apply(mv)
					val, ok := a.valueFn[b.String()]
					if !ok {
						a.logger.Printf("Couldn't find the position in the map!")
						return models.Move{X: 0, Y: 0, XO: a.turn}
					}
					if val > maxVal {
						moveToPlay = mv
						maxVal = val
					}
					b.RevertLastMove()
				}

				b.Apply(moveToPlay)
				strB := b.String()
				b.RevertLastMove()
				strA := b.String()

				a.valueFn[strA] += a.params.LearningRate * (a.valueFn[strB] - a.valueFn[strA])

				return moveToPlay
			}
		} else {
			if exploratory {
				moves := b.GetAllMoves()
				ind := rand.Intn(len(moves))
				return moves[ind]
			} else {
				moves := b.GetAllMoves()
				if len(moves) == 0 {
					a.logger.Println("Got 0 available legal moves on the board")
					return models.Move{X: 0, Y: 0, XO: a.turn}
				}

				moveToPlay := moves[0]
				b.Apply(moveToPlay)
				minVal, ok := a.valueFn[b.String()]
				if !ok {
					a.logger.Printf("Couldn't find the position in the map!")
					return models.Move{X: 0, Y: 0, XO: a.turn}
				}
				b.RevertLastMove()

				for _, mv := range moves {
					b.Apply(mv)
					val, ok := a.valueFn[b.String()]
					if !ok {
						a.logger.Printf("Couldn't find the position in the map!")
						return models.Move{X: 0, Y: 0, XO: a.turn}
					}
					if val < minVal {
						moveToPlay = mv
						minVal = val
					}
					b.RevertLastMove()
				}

				b.Apply(moveToPlay)
				strB := b.String()
				b.RevertLastMove()
				strA := b.String()

				a.valueFn[strA] += a.params.LearningRate * (a.valueFn[strB] - a.valueFn[strA])

				return moveToPlay
			}
		}
	}
}
