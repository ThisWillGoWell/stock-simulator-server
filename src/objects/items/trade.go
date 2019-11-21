package items

import (
	"github.com/ThisWillGoWell/stock-simulator-server/src/database"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/effect"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/notification"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/sender"
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
	item := Items[p.ParentItemUuid]

	newEffect, deleteEffect, err := effect.NewTradeEffect(p.TargetPortfolio, item.Name, item.Name, tradeEffect, p.Duration.Duration)
	if err != nil {
		return nil, err
	}

	notification.NotificationLock.Acquire("ActivateTradeEffectItem")
	defer notification.NotificationLock.Release()

	n := notification.NewEffectNotification(p.TargetPortfolio, newEffect.Title)
	n2 := notification.UsedItemNotification(p.TargetPortfolio, item.Uuid, item.Type)
	if dbErr := database.Db.Execute([]interface{}{newEffect, n, n2}, []interface{}{item, deleteEffect}); dbErr != nil {
		notification.DeleteNotification(n)
		notification.DeleteNotification(n2)
	}
	DeleteItem(item)

	effect.DeleteEffect(deleteEffect)
	wires.EffectsNewObject.Offer(newEffect.Effect)
	sender.SendDeleteObject(item.PortfolioUuid, item.Item)
	sender.SendDeleteObject(deleteEffect.PortfolioUuid, deleteEffect.Effect)
	sender.SendNewObject(n.PortfolioUuid, n.Notification)
	sender.SendNewObject(n2.PortfolioUuid, n2.Notification)
	sender.SendDeleteObject(item.Uuid, item.Item)
	return nil, nil
}
