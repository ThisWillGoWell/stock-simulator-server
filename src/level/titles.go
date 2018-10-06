package level

import (
	"github.com/stock-simulator-server/src/money"
)

var Levels = populateLevels()

type Level struct {
	Num            int64 `json:"num"`
	Cost           int64 `json:"cost"`
	MaxSharesStock int64 `json:"max_shares"`
}

func makeLevel(targetMap map[int64]*Level, level, cost, maxSharesStock int64) {
	targetMap[level] = &Level{
		Num:            level,
		Cost:           cost,
		MaxSharesStock: maxSharesStock,
	}
}

func populateLevels() map[int64]*Level {
	levels := make(map[int64]*Level)
	makeLevel(levels, 0, 0, 25)
	makeLevel(levels, 1, 2*money.Thousand, 25)
	makeLevel(levels, 2, 10*money.Thousand, 50)
	makeLevel(levels, 3, 50*money.Thousand, 50)
	makeLevel(levels, 4, 100*money.Thousand, 50)
	makeLevel(levels, 5, 100*money.Thousand, 50)
	makeLevel(levels, 6, 100*money.Thousand, 50)
	makeLevel(levels, 7, 100*money.Thousand, 50)
	makeLevel(levels, 8, 100*money.Thousand, 50)
	return levels
}
