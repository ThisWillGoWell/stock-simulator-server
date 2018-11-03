package effect

import (
	"time"
)

const TradeEffectType = "trade"

type TradingType struct {
}

func (TradingType) Name() string {
	return TradeEffectType
}

type TradeEffect struct {
	Uuid *string `json:"-"`

	FeeAmount     *int64   `json:"fee_amount,omitempty"`     // fee on all trades
	FeeMultiplier *float64 `json:"fee_multiplier,omitempty"` // fee % of the total fees

	ProfitAmount            *int64   `json:"profit_amount,omitempty"`     // profit amount %
	ProfitPercent           *float64 `json:"profit_percent,omitempty"`    // profit percent increase
	CurrentProfitMultiplier *float64 `json:"profit_multiplier,omitempty"` // current profit multiplier
	TradeBlocked            *bool    `json:"trade_blocked,omitempty"`     // if trade is blocked
}

func newTradingEffect(trade *TradeEffect, portfolioUuid, title string, duration time.Duration) {

}

func NewProfitEffect(portfolioUuid, title string, amount float64, duration time.Duration) {
	EffectLock.Acquire("new-protfit-effect")
	defer EffectLock.Release()
	newEffect := &TradeEffect{
		ProfitPercent: &amount,
	}
	newTradingEffect(newEffect, portfolioUuid, title, duration)
}

func NewBlockTrading(portfolioUuid, title string, duration time.Duration) {
	EffectLock.Acquire("new-block-effect")
	defer EffectLock.Release()
	newEffect := &TradeEffect{
		TradeBlocked: createBool(true),
	}
	newTradingEffect(newEffect, portfolioUuid, title, duration)
}

func FeelessTradeing(portfolioUuid, title string, duration time.Duration) {
	EffectLock.Acquire("new-block-effect")
	defer EffectLock.Release()
	newEffect := &TradeEffect{
		FeeMultiplier: createFloat(0),
	}
	newTradingEffect(newEffect, portfolioUuid, title, duration)
}

func (t TradeEffect) Add(effect *TradeEffect) {

	if effect.FeeAmount != nil {
		if t.FeeAmount == nil {
			t.FeeAmount = createInt(*effect.FeeAmount)
		} else {
			*t.FeeAmount += *effect.FeeAmount
		}
	}

	if effect.FeeMultiplier != nil {
		if t.FeeMultiplier == nil {
			t.FeeMultiplier = createFloat(*effect.FeeMultiplier)
		} else {
			*t.FeeMultiplier += *effect.FeeMultiplier
		}
	}

	if effect.ProfitAmount != nil {
		if t.ProfitAmount == nil {
			t.ProfitAmount = createInt(*effect.ProfitAmount)
		} else {
			*t.ProfitAmount += *effect.ProfitAmount
		}
	}

	if effect.ProfitPercent != nil {
		if t.ProfitPercent == nil {
			t.ProfitPercent = createFloat(*effect.ProfitPercent)
		} else {
			*t.ProfitPercent += *effect.ProfitPercent
		}
	}
	if effect.CurrentProfitMultiplier != nil {
		if t.CurrentProfitMultiplier == nil {
			t.CurrentProfitMultiplier = createFloat(*effect.CurrentProfitMultiplier)
		} else {
			*t.CurrentProfitMultiplier += *effect.CurrentProfitMultiplier
		}
	}

	if effect.TradeBlocked != nil {
		if *effect.TradeBlocked == true {
			t.TradeBlocked = createBool(true)
		}
	}
}

func createInt(x int64) *int64 {
	return &x
}
func createFloat(x float64) *float64 {
	return &x
}

func createBool(x bool) *bool {
	return &x
}
