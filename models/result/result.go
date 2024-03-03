package result

import (
	"tictactoe/models"
)

type Result struct {
	Draw   string
	Winner models.Turn
	Loser  models.Turn
}
