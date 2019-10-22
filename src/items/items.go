package items

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/merge"

	"github.com/pkg/errors"
	"github.com/ThisWillGoWell/stock-simulator-server/src/change"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/notification"
	"github.com/ThisWillGoWell/stock-simulator-server/src/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/sender"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"
)

var ItemsPortInventory = make(map[string]map[string]*Item)
var Items = make(map[string]*Item)
var ItemLock = lock.NewLock("item")

const ItemIdentifiableType = "item"

type InnerItem interface {
	SetPortfolioUuid(string)
	Activate(interface{}) (interface{}, error)
	Copy() InnerItem
	SetParentItemUuid(string)
}

type Item struct {
	Uuid          string           `json:"uuid"`
	Name          string           `json:"name"`
	ConfigId      string           `json:"config"`
	Type          string           `json:"type"`
	PortfolioUuid string           `json:"portfolio_uuid"`
	UpdateChannel chan interface{} `json:"-"`
	InnerItem     InnerItem        `json:"-" change:"inner"`
	CreateTime    time.Time        `json:"create_time"`
}

type i2 struct {
	Uuid          string    `json:"uuid"`
	Name          string    `json:"name"`
	ConfigId      string    `json:"config"`
	Type          string    `json:"type"`
	PortfolioUuid string    `json:"portfolio_uuid"`
	CreateTime    time.Time `json:"create_time"`
}

func (i *Item) GetId() string {
	return i.Uuid
}

func (*Item) GetType() string {
	return ItemIdentifiableType
}

func newItem(portfolioUuid, configId, itemType, name string, innerItem interface{}) *Item {

	item := MakeItem(utils.SerialUuid(), portfolioUuid, configId, itemType, name, innerItem, time.Now())
	sender.SendNewObject(portfolioUuid, item)
	wires.ItemsNewObjects.Offer(item)
	return item

}

func MakeItem(uuid, portfolioUuid, itemConfigId, itemType, name string, innerItem interface{}, createTime time.Time) *Item {
	switch innerItem.(type) {
	case string:
		innerItem = UnmarshalJsonItem(itemType, innerItem.(string))
	}
	i := &Item{
		Name:          name,
		ConfigId:      itemConfigId,
		Uuid:          uuid,
		PortfolioUuid: portfolioUuid,
		Type:          itemType,
		InnerItem:     innerItem.(InnerItem),
		UpdateChannel: make(chan interface{}),
		CreateTime:    createTime,
	}
	utils.RegisterUuid(uuid, i)

	if _, ok := ItemsPortInventory[i.PortfolioUuid]; !ok {
		ItemsPortInventory[i.PortfolioUuid] = make(map[string]*Item)
	}
	i.InnerItem.SetParentItemUuid(i.Uuid)
	ItemsPortInventory[i.PortfolioUuid][i.Uuid] = i
	Items[i.Uuid] = i
	change.RegisterPrivateChangeDetect(i, i.UpdateChannel)
	sender.RegisterChangeUpdate(i.PortfolioUuid, i.UpdateChannel)
	return i
}

func BuyItem(portUuid, configId string) (string, error) {

	port, _ := portfolio.GetPortfolio(portUuid)

	config, found := validItems[configId]
	if !found {
		return "", errors.New("config not found")
	}

	port.Lock.Acquire("buy item")
	defer port.Lock.Release()
	ItemLock.Acquire("buy-item")
	defer ItemLock.Release()

	if config.RequiredLevel > port.Level {
		return "", errors.New("not high enough level")
	}
	if config.Cost > port.Wallet {
		return "", errors.New("not enough $$ in wallet")
	}

	port.Wallet -= config.Cost
	if _, ok := ItemsPortInventory[port.Uuid]; !ok {
		ItemsPortInventory[port.Uuid] = make(map[string]*Item)
	}

	//todo

	newItem := newItem(portUuid, configId, config.Type, config.Name, config.Prams.Copy())
	newItem.InnerItem.SetPortfolioUuid(portUuid)
	ItemsPortInventory[port.Uuid][newItem.PortfolioUuid] = newItem
	Items[newItem.Uuid] = newItem

	notification.NewItemNotification(portUuid, newItem.Type, newItem.Uuid)
	go port.Update()
	return newItem.Uuid, nil
}

func (i *Item) DeleteItem() error {
	return DeleteItem(i.Uuid, i.PortfolioUuid, true)
}

func DeleteItem(uuid, portfolioUuid string, lockAcquired bool) error {
	if !lockAcquired {
		ItemLock.Acquire("delete-item")
		defer ItemLock.Release()
	}

	if _, exists := ItemsPortInventory[portfolioUuid]; !exists {
		return errors.New("user does not have any item")
	}
	item, exists := ItemsPortInventory[portfolioUuid][uuid]
	if !exists {
		return errors.New("item does not exist")
	}

	change.UnregisterChangeDetect(item)
	close(item.UpdateChannel)
	delete(Items, uuid)
	delete(ItemsPortInventory[item.PortfolioUuid], uuid)
	if len(ItemsPortInventory[item.PortfolioUuid]) == 0 {
		delete(ItemsPortInventory, item.PortfolioUuid)
	}
	utils.RemoveUuid(uuid)
	sender.SendDeleteObject(portfolioUuid, item)
	wires.ItemsDelete.Offer(item)
	return nil
}

func GetItemsForPortfolio(portfolioUuid string) []*Item {
	ItemLock.Acquire("get-Items")
	defer ItemLock.Release()
	items := make([]*Item, 0)
	userItems, ok := ItemsPortInventory[portfolioUuid]
	if !ok {
		return items
	}
	for _, item := range userItems {
		items = append(items, item)
	}
	return items
}

func getItem(itemId, portfolioUuid string) (*Item, error) {
	userItems, ok := ItemsPortInventory[portfolioUuid]
	if !ok {
		return nil, errors.New("user has no items")
	}
	item, ok := userItems[itemId]
	if !ok {
		return nil, errors.New("user does not have that item")
	}
	return item, nil
}

//func ViewItem(itemId, userUuid string) (interface{}, error) {
//	ItemLock.Acquire("Use Item")
//	defer ItemLock.Release()
//	item, err := getItem(itemId, userUuid)
//	if err != nil {
//		return nil, err
//	}
//	if !item.HasBeenUsed() {
//		return nil, errors.New("Item has not been used")
//	}
//	return item.View(), nil
//}

func Use(itemId, portfolioUuid string, itemParameters interface{}) (interface{}, error) {
	ItemLock.Acquire("Use Item")
	defer ItemLock.Release()
	item, err := getItem(itemId, portfolioUuid)
	if err != nil {
		return nil, err
	}

	val, err := item.InnerItem.Activate(itemParameters)
	if err != nil {
		notification.UsedItemNotification(portfolioUuid, itemId, item.Type)
	}
	return val, err
}

/**
 *
 */
func TransferItem(currentOwner, newOwner, itemId string) error {
	if _, ok := ItemsPortInventory[currentOwner]; !ok {
		return errors.New("current owner does not own any items")
	}
	item, ok := ItemsPortInventory[currentOwner][itemId]
	if !ok {
		return errors.New("current owner does not have the item id")
	}

	if _, ok := ItemsPortInventory[newOwner]; !ok {
		ItemsPortInventory[newOwner] = make(map[string]*Item)
	}
	ItemsPortInventory[currentOwner][itemId] = item

	delete(ItemsPortInventory[currentOwner], itemId)
	if len(ItemsPortInventory[currentOwner]) == 0 {
		delete(ItemsPortInventory, currentOwner)
	}
	return nil
}

func UnmarshalJsonItem(itemType, jsonStr string) InnerItem {
	var item InnerItem
	switch itemType {
	case TradeItemType:
		item = &TradeEffectItem{}
	}
	err := json.Unmarshal([]byte(jsonStr), &item)
	if err != nil {
		log.Fatal("error unmarshal json item", err.Error())
	}
	return item
}

func (i *Item) MarshalJSON() ([]byte, error) {

	return merge.Json(i2{
		PortfolioUuid: i.PortfolioUuid,
		Uuid:          i.Uuid,
		Name:          i.Name,
		ConfigId:      i.ConfigId,
		CreateTime:    i.CreateTime,
		Type:          i.Type,
	}, i.InnerItem)
}
