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
	makeLevel(levels, 1, 10*money.Thousand, 25)
	makeLevel(levels, 2, 250*money.Thousand, 25)
	makeLevel(levels, 3, 500*money.Thousand, 50)
	makeLevel(levels, 4, 1*money.Million, 50)
	makeLevel(levels, 5, 1500*money.Thousand, 75)
	makeLevel(levels, 6, 2500*money.Thousand, 75)
	makeLevel(levels, 7, 5*money.Million, 75)
	makeLevel(levels, 8, 10*money.Million, 100)
	return levels
}
