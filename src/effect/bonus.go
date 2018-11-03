package effect

import (
	"time"

	"github.com/stock-simulator-server/src/utils"
)

type BonusTradingType struct {
}

func (BonusTradingType) Name() string {
	return "bonus_trading"
}
func (BonusTradingType) Cost() int64 {
	return 0
}

func (BonusTradingType) RequiredLevel() int64 {
	return 0
}

type BonusTrading struct {
	BonusAmount float64        `json:"bonus"`
}

func (bonus *BonusTrading) IsPermanent() bool {
	return true
}

func NewBonusTrading(portfolioUuid, title string, amount float64, duration utils.Duration) error{
	EffectLock.Acquire("new-bonus-trading")
	defer EffectLock.Release()

	newEffect := &Effect{
		Uuid: utils.SerialUuid(),
		Title: title,
		Active: true,
		StartTime: time.Now(),
		Type: BonusTradingType{},
		InnerEffect: &BonusTrading{
			BonusAmount: amount,
		},
	}

	effects, ok := PortfolioEffects[portfolioUuid]
	if !ok{
		n
	}
	return nil
}

func