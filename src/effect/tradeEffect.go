package effect

import (
	"time"

	"github.com/stock-simulator-server/src/wires"

	"github.com/stock-simulator-server/src/money"
)

const TradeEffectType = "trade"
const baseTradeEffectTag = "base_trade"

const BaseTaxRate = 0.15
const BaseSellFell = 20 * money.One
const BaseBuyFee = 20 * money.One

type TradingType struct {
}

func (TradingType) Name() string {
	return TradeEffectType
}

type TradeEffect struct {
	parentEffect *Effect `json:"-"`

	BuyFeeAmount     *int64   `json:"buy_fee_amount,omitempty" change:"-"`     // fee on all trades ex: base fee
	BuyFeeMultiplier *float64 `json:"buy_fee_multiplier,omitempty" change:"-"` // fee % of the total fees, ex: double fees

	SellFeeAmount     *int64   `json:"sell_fee_amount,omitempty" change:"-"`     // fee on all sales trades ex: base fee
	SellFeeMultiplier *float64 `json:"sell_fee_multiplier,omitempty" change:"-"` // fee % of the total fees, ex: double fees on trades

	BonusProfitMultiplier *float64 `json:"profit_multiplier,omitempty" change:"-"` // current profit multiplier, ex: bonus

	TaxPercent    *float64 `json:"tax_percent,omitempty" change:"-"`    // tax payed on profits, ex:  base tax
	TaxMultiplier *float64 `json:"tax_multiplier,omitempty" change:"-"` // percent multiplier, ex: taxless sales

	TradeBlocked *bool `json:"trade_blocked,omitempty" change:"-"` // if trade is blocked
}

func NewBaseTradeEffect(portfolioUuid string) {
	baseTradeEffect := &TradeEffect{
		BuyFeeAmount:  createInt(BaseBuyFee),
		SellFeeAmount: createInt(BaseSellFell),
		TaxPercent:    createFloat(BaseTaxRate),
	}
	baseTradeEffect.parentEffect = newEffect(portfolioUuid, "Base Effect", baseTradeEffectTag, TradeEffectType, baseTradeEffect, 0)
}

func UpdateBaseProfit(portfolioUuid string, profitMultiplier float64) {
	EffectLock.Acquire("update-effect")
	defer EffectLock.Release()
	effect := getTaggedEffect(portfolioUuid, baseTradeEffectTag)
	effect.InnerEffect.(*TradeEffect).BonusProfitMultiplier = createFloat(profitMultiplier)
	wires.EffectsUpdate.Offer(effect)

}

func NewProfitEffect(portfolioUuid, title string, amount float64, duration time.Duration) {
	newTradeEffect := &TradeEffect{
		BonusProfitMultiplier: &amount,
	}
	newTradeEffect.parentEffect = newEffect(portfolioUuid, title, TradeEffectType, "", newTradeEffect, duration)

}

func NewBlockTrading(portfolioUuid, title string, duration time.Duration) {
	newTradeEffect := &TradeEffect{
		TradeBlocked: createBool(true),
	}
	newTradeEffect.parentEffect = newEffect(portfolioUuid, title, TradeEffectType, "", newTradeEffect, duration)
}

func NewFeelessTradeing(portfolioUuid, title string, duration time.Duration) {
	newTradeEffect := &TradeEffect{
		SellFeeMultiplier: createFloat(0),
		BuyFeeMultiplier:  createFloat(0),
	}
	newTradeEffect.parentEffect = newEffect(portfolioUuid, title, "", TradeEffectType, newTradeEffect, duration)
}

func NewTaxModifier(portfolioUuid, title string, duration time.Duration, taxMultiplier float64) {
	newTradeEffect := &TradeEffect{
		TaxMultiplier: &taxMultiplier,
	}
	newTradeEffect.parentEffect = newEffect(portfolioUuid, title, "", TradeEffectType, newTradeEffect, duration)

}

// Calculate the total bonus  for a portfolio
func TotalTradeEffect(portfolioUuid string) (*TradeEffect, []string) {
	EffectLock.Acquire("TotalBonus")
	defer EffectLock.Release()
	effect, ok := portfolioEffects[portfolioUuid]
	totalEffect := &TradeEffect{
		BuyFeeAmount:     createInt(0),
		BuyFeeMultiplier: createFloat(1),

		SellFeeAmount:     createInt(0),
		SellFeeMultiplier: createFloat(1),

		BonusProfitMultiplier: createFloat(0),

		TaxPercent:    createFloat(0),
		TaxMultiplier: createFloat(1),

		TradeBlocked: createBool(false),
	}
	uuids := make([]string, 0)
	if !ok {
		return totalEffect, uuids
	}
	for uuid, e := range effect {
		switch e.Type {
		case TradeEffectType:
			totalEffect.Add(e.InnerEffect.(*TradeEffect))
			uuids = append(uuids, uuid)
		}
	}
	return totalEffect, uuids
}

func (t *TradeEffect) Add(effect *TradeEffect) {

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

	if effect.BonusProfitMultiplier != nil {
		if t.BonusProfitMultiplier == nil {
			t.BonusProfitMultiplier = createFloat(*effect.BonusProfitMultiplier)
		} else {
			*t.BonusProfitMultiplier += *effect.BonusProfitMultiplier
		}
	}
	if effect.SellFeeMultiplier != nil {
		if t.SellFeeMultiplier == nil {
			t.SellFeeMultiplier = createFloat(*effect.SellFeeMultiplier)
		} else {
			*t.SellFeeMultiplier += *effect.SellFeeMultiplier
		}
	}
	if effect.TaxMultiplier != nil {
		if t.TaxMultiplier == nil {
			t.TaxMultiplier = createFloat(*effect.TaxMultiplier)
		} else {
			*t.TaxMultiplier += *effect.TaxMultiplier
		}
	}
	if effect.TaxPercent != nil {
		if t.TaxPercent == nil {
			t.TaxPercent = createFloat(*effect.TaxPercent)
		} else {
			*t.TaxPercent += *effect.TaxPercent
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

func (t *TradeEffect) GetId() string {
	return t.parentEffect.Uuid
}

func (*TradeEffect) GetType() string {
	return EffectIdType
}
