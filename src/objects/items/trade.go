package items

import (
	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/effect"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
)

const TradeItemType = "trade_effect"

type TradeEffectItem struct {
	ParentItemUuid    string         `json:"-"`
	TargetPortfolio   string         `json:"portfolio_uuid"`
	Duration          utils.Duration `json:"duration"`
	ProfitMultiplier  *float64       `json:"profit_multiplier"`
	BuyFeeMultiplier  *float64       `json:"buy_fee_multiplier"`
	SellFeeMultiplier *float64       `json:"sell_fee_multiplier"`
	BuyFeeAmount      *int64         `json:"buy_fee"`
	SellFeeAmount     *int64         `json:"sell_fee"`
	BlockTrading      *bool          `json:"block_trading"`
}

func (p *TradeEffectItem) SetPortfolioUuid(portfolioUuid string) {
	p.TargetPortfolio = portfolioUuid
}

func (p *TradeEffectItem) SetParentItemUuid(parent string) {
	p.ParentItemUuid = parent
}

func (p *TradeEffectItem) Copy() InnerItem {
	return &TradeEffectItem{
		Duration:          p.Duration,
		BuyFeeAmount:      p.BuyFeeAmount,
		BuyFeeMultiplier:  p.BuyFeeMultiplier,
		SellFeeMultiplier: p.SellFeeMultiplier,
		SellFeeAmount:     p.SellFeeAmount,
		ProfitMultiplier:  p.ProfitMultiplier,
	}
}

func (p *TradeEffectItem) Activate(interface{}) (interface{}, error) {

	tradeEffect := &effect.TradeEffect{
		BuyFeeAmount:          p.BuyFeeAmount,
		BuyFeeMultiplier:      p.BuyFeeMultiplier,
		SellFeeMultiplier:     p.SellFeeMultiplier,
		SellFeeAmount:         p.SellFeeAmount,
		BonusProfitMultiplier: p.ProfitMultiplier,
	}
	parent := Items[p.ParentItemUuid]
	if err := effect.NewTradeEffect(p.TargetPortfolio, parent.Name, parent.Name, tradeEffect, p.Duration.Duration); err != nil {
		return nil, err
	}
	if err := parent.DeleteItem(); err != nil {
		log.Log.Error("failed to delete item (in database) err=[%v]", err)
	}
	return nil, nil
}
