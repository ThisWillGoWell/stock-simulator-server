package level

import (
	"encoding/json"

	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"
	"github.com/ThisWillGoWell/stock-simulator-server/src/game/money"
)

var Levels = make(map[int64]*Level)

type Level struct {
	Num              int64   `json:"num"`
	Cost             int64   `json:"cost"`
	MaxSharesStock   int64   `json:"max_shares"`
	ProfitMultiplier float64 `json:"profit_multiplier"`
}

func makeLevel(targetMap map[int64]*Level, level, cost, maxSharesStock int64, profitMultiplier float64) {
	targetMap[level] = &Level{
		Num:              level,
		Cost:             cost,
		MaxSharesStock:   maxSharesStock,
		ProfitMultiplier: profitMultiplier,
	}
}

func LoadLevels(data []byte) {
	levelsList := make([]*Level, 0)
	err := json.Unmarshal(data, &levelsList)
	if err != nil {
		log.Log.Error("err loading levels", err)
	}
	for i, ele := range levelsList {
		Levels[int64(i)] = ele
	}
}

func populateLevels() map[int64]*Level {
	levels := make(map[int64]*Level)
	makeLevel(levels, 0, 0, 25, 0)
	makeLevel(levels, 1, 10*money.Thousand, 25, 0.02)
	makeLevel(levels, 2, 250*money.Thousand, 25, 0.04)
	makeLevel(levels, 3, 500*money.Thousand, 50, 0.07)
	makeLevel(levels, 4, 1*money.Million, 50, 0.1)
	makeLevel(levels, 5, 1500*money.Thousand, 75, 0.12)
	makeLevel(levels, 6, 2500*money.Thousand, 75, 0.14)
	makeLevel(levels, 7, 5*money.Million, 75, 0.18)
	makeLevel(levels, 8, 10*money.Million, 100, 0.2)
	return levels
}
