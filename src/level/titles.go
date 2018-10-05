package level

import (
	"github.com/stock-simulator-server/src/money"
)

var Levels = populateLevels()

type Level struct {
	Num  int64 `json:"num"`
	Cost int64 `json:"cost"`
}

func makeLevel(targetMap map[int64]*Level, level, cost int64) {
	targetMap[level] = &Level{
		Num:  level,
		Cost: cost,
	}
}

func populateLevels() map[int64]*Level {
	levels := make(map[int64]*Level)
	makeLevel(levels, 0, 0)
	makeLevel(levels, 1, 2*money.Thousand)
	makeLevel(levels, 2, 10*money.Thousand)
	makeLevel(levels, 3, 50*money.Thousand)
	makeLevel(levels, 3, 100*money.Thousand)
	return levels
}
