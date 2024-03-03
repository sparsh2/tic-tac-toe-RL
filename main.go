package main

import (
	"fmt"
	"log"
	"tictactoe/models"
	"tictactoe/models/board"
	"tictactoe/models/player"
)

func getValueDefault(board *board.Board, mp map[string]float64, c *int) {
	done := board.IsDone()
	decisive := board.IsDecisive()

	mp[board.String()] = 0.5

	if done && !decisive {
		*c = *c + 1
		mp[board.String()] = 0.5
	}

	if decisive {
		*c += 1
		res := board.GetResult()
		if res.Winner == models.X {
			mp[board.String()] = 1.0
		} else {
			mp[board.String()] = 0.0
		}
	}

	moves := board.GetAllMoves()
	for _, mv := range moves {
		board.Apply(mv)
		getValueDefault(board, mp, c)
		board.RevertLastMove()
	}
}

func main() {
	board := &board.Board{}
	board.Reset()

	valFn := map[string]float64{}

	c := 0

	getValueDefault(board, valFn, &c)

	logger := log.Default()
	agentLogger := log.Default()
	agentLogger.SetPrefix("Agent")
	agentX := player.NewAgent(valFn, &player.AgentParameters{
		Mode:                   player.LearnMode,
		ExploratoryProbability: 0.1,
		LearningRate:           0.01,
	}, agentLogger, models.X)

	agentO := player.NewAgent(valFn, &player.AgentParameters{
		Mode:                   player.LearnMode,
		ExploratoryProbability: 0.1,
		LearningRate:           0.1,
	}, agentLogger, models.O)

	times := 10000

	for i := 0; i < times; i++ {
		board.Reset()

		for !board.IsDone() && !board.IsDecisive() {
			turn := board.GetTurn()

			if turn == models.X {
				mv := agentX.Move(*board)
				ok := board.Check(mv)
				if !ok {
					logger.Println("agentX gave an illegal move")
					break
				}
				board.Apply(mv)
			} else {
				mv := agentO.Move(*board)
				ok := board.Check(mv)
				if !ok {
					logger.Println("agentO gave an illegal move")
					break
				}
				board.Apply(mv)
			}
		}

		// board.PrintHistory()

		if i%100 == 0 || i == times-1 {
			board.PrintHistory()
		}
	}

	times = 20
	fmt.Println(c)
	agentX = player.NewAgent(valFn, &player.AgentParameters{
		Mode:                   player.PlayMode,
		ExploratoryProbability: 0.0,
		LearningRate:           0.01,
	}, agentLogger, models.X)
	agentO = player.NewAgent(valFn, &player.AgentParameters{
		Mode:                   player.PlayMode,
		ExploratoryProbability: 0.0,
		LearningRate:           0.01,
	}, agentLogger, models.O)
	for i := 0; i < times; i++ {
		board.Reset()

		for !board.IsDone() && !board.IsDecisive() {
			turn := board.GetTurn()

			if turn == models.X {
				mv := agentX.Move(*board)
				ok := board.Check(mv)
				if !ok {
					logger.Println("agentX gave an illegal move")
					break
				}
				board.Apply(mv)
				board.Print()
				fmt.Printf("value of board: %f\n", valFn[board.String()])
			} else {
				mv := getUserMove(board)
				board.Apply(mv)
				board.Print()
				fmt.Printf("value of board: %f\n", valFn[board.String()])
			}
		}

		res := board.GetResult()
		if res.Draw == "draw" {
			fmt.Printf("Game drawn\n")
		} else {
			fmt.Printf("Winner is %s\n", res.Winner.String())
		}
	}
}

func getUserMove(b *board.Board) models.Move {
	var x int
	var y int

	fmt.Scanf("%d", &x)
	fmt.Scanf("%d", &y)

	for !b.Check(models.Move{X: x, Y: y, XO: models.O}) {
		fmt.Scanf("%d", &x)
		fmt.Scanf("%d", &y)
	}

	return models.Move{X: x, Y: y, XO: models.O}
}
