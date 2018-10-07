package effect

import (
	"time"

	"github.com/stock-simulator-server/src/utils"
)

type BonusTradingType struct {
}

func (*BonusTradingType) Name() string {
	return "bonus_trading"
}
func (*BonusTradingType) Cost() int64 {
	return 0
}

func (*BonusTradingType) RequiredLevel() int64 {
	return 0
}

type BonusTrading struct {
	Duration  utils.Duration `json:"duration"`
	StartTime time.Time      `json:"start_time"`
}

func (bonus *BonusTrading) IsPermanent() bool {
	return true
}

func (bonus *BonusTrading) StartTime() int64 {
	return bonus.StartTime
}

func (*BonusTrading) GetType() EffectType {
	return &BonusTradingType{}
}
