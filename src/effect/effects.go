package effect

import (
	"time"

	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/utils"
)

var EffectLock = lock.NewLock("EffectLock")
var PortfolioEffects = make(map[string]map[string]*Effect)

// Calculate the total bonus  for a portfolio
func TotalBonus(portfolioUuid string) float64 {
	EffectLock.Acquire("TotalBonus")
	effect, ok := PortfolioEffects[portfolioUuid]
	if !ok {
		return 0
	}
	totalBonus := 0.0
	for _, e := range effect {
		switch e.InnerEffect.(type) {
		case BonusTradingType:
			totalBonus += e.InnerEffect.(*BonusTrading).BonusAmount
		}
	}
	return totalBonus
}

type EffectType interface {
	Name() string
	Cost() int64
	RequiredLevel() int64
}

type Effect struct {
	PortfolioUuid string         `json:"portfolio_uuid"`
	Uuid          string         `json:"uuid"`
	Title         string         `json:"title"`
	Active        bool           `json:"active"`
	IsPublic      bool           `json:"public"`
	Duration      utils.Duration `json:"duration"`
	StartTime     time.Time      `json:"time"`
	Type          EffectType     `json:"effect"`
	InnerEffect   interface{}
}
