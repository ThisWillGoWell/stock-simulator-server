package items

import (
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

type InsiderTradingItem struct {
	Type     InsiderTraderItemType
	UserUuid string           `json:"user_uuid"`
	Uuid     string           `json:"uuid"`
	Used     bool             `json:"used"`
	Result   map[string]int64 `json:"view,omitempty"`
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

func (it *InsiderTradingItem) GetUserUuid() string {
	return it.UserUuid
}
func (it *InsiderTradingItem) GetUuid() string {
	return it.Uuid
}
func (it *InsiderTradingItem) SetUserUuid(uuid string) {
	it.UserUuid = uuid
}

func (it *InsiderTradingItem) HasBeenUsed() bool {
	return it.Used
}

func newInsiderTradingItem(userUuid string) *InsiderTradingItem {
	uuid := utils.SerialUuid()
	item := &InsiderTradingItem{
		UserUuid: userUuid,
		Uuid:     uuid,
		Used:     false,
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
