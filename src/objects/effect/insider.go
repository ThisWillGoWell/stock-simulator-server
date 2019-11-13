package effect

import (
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
)

const InsiderTradingEffectType = "insider"

type InsiderTradingType struct {
}

func (*InsiderTradingType) Name() string {
	return InsiderTradingEffectType
}

func (*InsiderTradingType) RequiredLevel() int64 {
	return 0
}

type InsiderTrading struct {
	Duration  utils.Duration
	StartTime time.Time
}

func (*InsiderTrading) Type() EffectType {
	return &InsiderTradingType{}
}
