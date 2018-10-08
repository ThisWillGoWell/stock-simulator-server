package items

import (
	"fmt"

	"github.com/stock-simulator-server/src/money"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/valuable"
)

const insiderTradingItemType = "insider"

type InsiderTraderItemType struct {
}
type InsiderTraderItemParameters struct {
}

func (InsiderTraderItemType) GetName() string {
	return "Insider Trading"
}

func (InsiderTraderItemType) GetCost() int64 {
	return 1 * money.Thousand
}

func (InsiderTraderItemType) GetType() string {
	return insiderTradingItemType
}

func (InsiderTraderItemType) GetDescription() string {
	return "View the current target prices of all the stocks you own"
}

func (InsiderTraderItemType) GetActivateParameters() interface{} {
	return InsiderTraderItemParameters{}
}

func (InsiderTraderItemType) RequiredLevel() int64 {
	return 1
}

func (InsiderTraderItemType) UnmarshalJSON(data []byte) error {
	return nil
}

type InsiderTradingItem struct {
	Type          InsiderTraderItemType
	PortfolioUuid string           `json:"portfolio_uuid"`
	Uuid          string           `json:"uuid"`
	Used          bool             `json:"used" change:"-"`
	Result        map[string]int64 `json:"target_prices,omitempty" change:"-"`
	UpdateChan    chan interface{} `json:"-"`
}

func (it *InsiderTradingItem) GetItemType() ItemType {
	return it.Type
}
func (it *InsiderTradingItem) GetType() string {
	return ItemIdentifiableType
}

func (it *InsiderTradingItem) GetId() string {
	return it.Uuid
}

func (it *InsiderTradingItem) GetPortfolioUuid() string {
	return it.PortfolioUuid
}
func (it *InsiderTradingItem) GetUuid() string {
	return it.Uuid
}
func (it *InsiderTradingItem) SetPortfolioUuid(uuid string) {
	it.PortfolioUuid = uuid
}

func (it *InsiderTradingItem) HasBeenUsed() bool {
	return it.Used
}

func (it *InsiderTradingItem) GetUpdateChan() chan interface{} {
	return it.UpdateChan
}

func (it *InsiderTradingItem) Load() {

}

func (it *InsiderTradingItem) GetItem() interface{} {
	return struct{}{}
}

func newInsiderTradingItem(userUuid string) *InsiderTradingItem {
	uuid := utils.SerialUuid()
	item := &InsiderTradingItem{
		PortfolioUuid: userUuid,
		Uuid:          uuid,
		Used:          false,
		UpdateChan:    make(chan interface{}),
	}
	utils.RegisterUuid(uuid, item)
	return item
}

func (it *InsiderTradingItem) Activate(interface{}) (interface{}, error) {
	valuable.ValuablesLock.Acquire("active insider trading")
	defer valuable.ValuablesLock.Release()

	targetPrices := make(map[string]int64)

	for _, stock := range valuable.Stocks {
		targetPrices[stock.Uuid] = stock.PriceChanger.GetTargetPrice()
	}
	it.Used = true
	it.Result = targetPrices
	return targetPrices, nil
}

func (it *InsiderTradingItem) View() interface{} {
	return it.Result
}

func (u *InsiderTraderItemType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, insiderTradingItemType)), nil
}
