package items

import (
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/money"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/valuable"
)

type InsiderTraderItemType struct {

}
type InsiderTraderItemParameters struct {

}

func (InsiderTraderItemType) GetName() string{
	return "Insider Trading"
}

func (InsiderTraderItemType) GetCost() int64{
	return 1 * money.Thousand
}

func (InsiderTraderItemType) GetDescription() string{
	return "View the current target prices of all the stocks you own"
}

func (InsiderTraderItemType) GetActivateParameters() interface{}{
	return InsiderTraderItemParameters{}
}

func (InsiderTraderItemType) RequiredLevel() int64{
	return 1
}

type InsiderTradingItem struct {
	Type InsiderTraderItemType
	UserUuid string
	Uuid string
	Used bool
}

func (it *InsiderTradingItem) GetType() ItemType{
	return it.Type
}
func  (it *InsiderTradingItem) GetUserUuid() string {
	return it.UserUuid
}
func  (it *InsiderTradingItem) GetUuid() string {
	return it.Uuid
}

func (it *InsiderTradingItem) HasBeenUsed() bool{
	return it.Used
}

func newInsiderTradingItem(userUuid string) *InsiderTradingItem{
	uuid := utils.PseudoUuid()
	item := &InsiderTradingItem{
		UserUuid: userUuid,
		Uuid: uuid,
		Used: false,
	}
	utils.RegisterUuid(uuid, item)
	return item
}


func  (it *InsiderTradingItem) Activate() (interface{}, error){
	portfolioUuid := account.UserList[it.Uuid].PortfolioId
	ledger.EntriesLock.Acquire("activate insider Trading")
	defer ledger.EntriesLock.Release()

	targetPrices := make(map[string]int64)
	ledgers, hasLedgers := ledger.EntriesPortfolioStock[portfolioUuid]
	if !hasLedgers {
		return targetPrices, nil
	}

	for _, l := range ledgers{
		targetPrices[l.StockId] = valuable.Stocks[l.StockId].PriceChanger.GetTargetPrice()
	}
	it.Used = true

	return targetPrices, nil
}