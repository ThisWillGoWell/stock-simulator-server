package effect

import (
	"time"

	"github.com/stock-simulator-server/src/utils"

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
		BuyFeeAmount:  utils.CreateInt(BaseBuyFee),
		SellFeeAmount: utils.CreateInt(BaseSellFell),
		TaxPercent:    utils.CreateFloat(BaseTaxRate),
	}
	baseTradeEffect.parentEffect = newEffect(portfolioUuid, "Base Effect", baseTradeEffectTag, TradeEffectType, baseTradeEffect, 0)
}

func UpdateBaseProfit(portfolioUuid string, profitMultiplier float64) {
	EffectLock.Acquire("update-effect")
	defer EffectLock.Release()
	effect := getTaggedEffect(portfolioUuid, baseTradeEffectTag)
	effect.InnerEffect.(*TradeEffect).BonusProfitMultiplier = utils.CreateFloat(profitMultiplier)
	wires.EffectsUpdate.Offer(effect)

}

func NewTradeEffect(portfolioUuid, title, tag string, effect *TradeEffect, duration time.Duration) {

	effect.parentEffect = newEffect(portfolioUuid, title, TradeEffectType, tag, effect, duration)
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
		BuyFeeAmount:     utils.CreateInt(0),
		BuyFeeMultiplier: utils.CreateFloat(1),

		SellFeeAmount:     utils.CreateInt(0),
		SellFeeMultiplier: utils.CreateFloat(1),

		BonusProfitMultiplier: utils.CreateFloat(0),

		TaxPercent:    utils.CreateFloat(0),
		TaxMultiplier: utils.CreateFloat(1),

		TradeBlocked: utils.CreateBool(false),
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
			t.BuyFeeAmount = utils.CreateInt(*effect.BuyFeeAmount)
		} else {
			*t.BuyFeeAmount += *effect.BuyFeeAmount
		}
	}

	if effect.BuyFeeMultiplier != nil {
		if t.BuyFeeMultiplier == nil {
			t.BuyFeeMultiplier = utils.CreateFloat(*effect.BuyFeeMultiplier)
		} else {
			*t.BuyFeeMultiplier += *effect.BuyFeeMultiplier
		}
	}

	if effect.SellFeeAmount != nil {
		if t.SellFeeAmount == nil {
			t.SellFeeAmount = utils.CreateInt(*effect.SellFeeAmount)
		} else {
			*t.SellFeeAmount += *effect.SellFeeAmount
		}
	}

	if effect.BonusProfitMultiplier != nil {
		if t.BonusProfitMultiplier == nil {
			t.BonusProfitMultiplier = utils.CreateFloat(*effect.BonusProfitMultiplier)
		} else {
			*t.BonusProfitMultiplier += *effect.BonusProfitMultiplier
		}
	}
	if effect.SellFeeMultiplier != nil {
		if t.SellFeeMultiplier == nil {
			t.SellFeeMultiplier = utils.CreateFloat(*effect.SellFeeMultiplier)
		} else {
			*t.SellFeeMultiplier += *effect.SellFeeMultiplier
		}
	}
	if effect.TaxMultiplier != nil {
		if t.TaxMultiplier == nil {
			t.TaxMultiplier = utils.CreateFloat(*effect.TaxMultiplier)
		} else {
			*t.TaxMultiplier += *effect.TaxMultiplier
		}
	}
	if effect.TaxPercent != nil {
		if t.TaxPercent == nil {
			t.TaxPercent = utils.CreateFloat(*effect.TaxPercent)
		} else {
			*t.TaxPercent += *effect.TaxPercent
		}
	}

	if effect.TradeBlocked != nil {
		if *effect.TradeBlocked == true {
			t.TradeBlocked = utils.CreateBool(true)
		}
	}
}

func (t *TradeEffect) GetId() string {
	return t.parentEffect.Uuid
}

func (*TradeEffect) GetType() string {
	return EffectIdType
}
