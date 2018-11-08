package effect

import (
	"time"
)

const TradeEffectType = "trade"
const BaseEffectTag = "tag"

type TradingType struct {
}

func (TradingType) Name() string {
	return TradeEffectType
}

type TradeEffect struct {
	Uuid *string `json:"-"`

	BuyFeeAmount     *int64   `json:"buy_fee_amount,omitempty"`     // fee on all trades
	BuyFeeMultiplier *float64 `json:"buy_fee_multiplier,omitempty"` // fee % of the total fees

	SellFeeAmount     *int64   `json:"sell_fee_amount,omitempty"`     // fee on all sales trades
	SellFeeMultiplier *float64 `json:"sell_fee_multiplier,omitempty"` // fee % of the total fees

	ProfitPercent    *float64 `json:"profit_percent,omitempty"`    // profit percent change
	ProfitMultiplier *float64 `json:"profit_multiplier,omitempty"` // current profit multiplier

	TaxPercent    *int64   `json:"tax_percent,omitempty"`
	TaxMultiplier *float64 `json:"tax_multiplier,omitempty"`

	TradeBlocked *bool `json:"trade_blocked,omitempty"` // if trade is blocked
}


func NewProfitEffect(portfolioUuid, title string, amount float64, duration time.Duration) {
	EffectLock.Acquire("new-protfit-effect")
	defer EffectLock.Release()
	newTradeEffect := &TradeEffect{
		ProfitPercent: &amount,
	}
	newEffect(portfolioUuid, title, TradeEffectType, newTradeEffect, duration)
}

func NewBlockTrading(portfolioUuid, title string, duration time.Duration) {
	EffectLock.Acquire("new-block-effect")
	defer EffectLock.Release()
	newTradeEffect := &TradeEffect{
		TradeBlocked: createBool(true),
	}
	newEffect(portfolioUuid, title, TradeEffectType, newTradeEffect, duration)
}

func FeelessTradeing(portfolioUuid, title string, duration time.Duration) {
	EffectLock.Acquire("new-block-effect")
	defer EffectLock.Release()
	newTradeEffect := &TradeEffect{
		SellFeeMultiplier: createFloat(0),
		BuyFeeMultiplier:  createFloat(0),
	}
	newEffect(portfolioUuid, title, TradeEffectType, newTradeEffect, duration)
}

func TaxReduction(portfolioUuid, title string, duration time.Duration, taxReductionAmount float64) {
	EffectLock.Acquire("new-block-effect")
	defer EffectLock.Release()
	newTradeEffect := &TradeEffect{
		TaxMultiplier: createFloat(taxReductionAmount),
	}
	newEffect(portfolioUuid, title, TradeEffectType, newTradeEffect, duration)
}

func (t TradeEffect) Add(effect *TradeEffect) {

	if effect.BuyFeeAmount != nil {
		if t.BuyFeeAmount == nil {
			t.BuyFeeAmount = createInt(*effect.BuyFeeAmount)
		} else {
			*t.BuyFeeAmount += *effect.BuyFeeAmount
		}
	}

	if effect.BuyFeeMultiplier != nil {
		if t.BuyFeeMultiplier == nil {
			t.BuyFeeMultiplier = createFloat(*effect.BuyFeeMultiplier)
		} else {
			*t.BuyFeeMultiplier += *effect.BuyFeeMultiplier
		}
	}

	if effect.SellFeeAmount != nil {
		if t.SellFeeAmount == nil {
			t.SellFeeAmount = createInt(*effect.SellFeeAmount)
		} else {
			*t.SellFeeAmount += *effect.SellFeeAmount
		}
	}

	if effect.ProfitPercent != nil {
		if t.ProfitPercent == nil {
			t.ProfitPercent = createFloat(*effect.ProfitPercent)
		} else {
			*t.ProfitPercent += *effect.ProfitPercent
		}
	}
	if effect.SellFeeMultiplier != nil {
		if t.SellFeeMultiplier == nil {
			t.SellFeeMultiplier = createFloat(*effect.SellFeeMultiplier)
		} else {
			*t.SellFeeMultiplier += *effect.SellFeeMultiplier
		}
	}

	if effect.TradeBlocked != nil {
		if *effect.TradeBlocked == true {
			t.TradeBlocked = createBool(true)
		}
	}
}

type BaseTradeEffect(){

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


