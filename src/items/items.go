package items

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/database"

	"github.com/ThisWillGoWell/stock-simulator-server/src/merge"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"github.com/ThisWillGoWell/stock-simulator-server/src/change"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/log"
	"github.com/ThisWillGoWell/stock-simulator-server/src/notification"
	"github.com/ThisWillGoWell/stock-simulator-server/src/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/sender"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
	"github.com/pkg/errors"
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

func newItem(portfolioUuid, configId, itemType, name string, innerItem interface{}) (i *Item, err error) {

	if i, err = MakeItem(utils.SerialUuid(), portfolioUuid, configId, itemType, name, innerItem, time.Now()); err != nil {
		log.Log.Errorf("making item err=[%v]", err)
		return nil, err
	}

	if err := database.Db.WriteItem(i); err != nil {
		_ = DeleteItem(i.Uuid, portfolioUuid, false, false, true)
		return nil, err
	}

	sender.SendNewObject(portfolioUuid, i)
	wires.ItemsNewObjects.Offer(i)
	return
}

func MakeItem(uuid, portfolioUuid, itemConfigId, itemType, name string, innerItem interface{}, createTime time.Time) (*Item, error) {
	switch innerItem.(type) {
	case string:
		var err error
		if innerItem, err = UnmarshalJsonItem(itemType, innerItem.(string)); err != nil {
			return nil, err
		}
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

	if _, ok := ItemsPortInventory[i.PortfolioUuid]; !ok {
		ItemsPortInventory[i.PortfolioUuid] = make(map[string]*Item)
	}
	i.InnerItem.SetParentItemUuid(i.Uuid)

	if err := change.RegisterPrivateChangeDetect(i, i.UpdateChannel); err != nil {
		return nil, err
	}

	utils.RegisterUuid(uuid, i)
	ItemsPortInventory[i.PortfolioUuid][i.Uuid] = i
	Items[i.Uuid] = i

	sender.RegisterChangeUpdate(i.PortfolioUuid, i.UpdateChannel)
	return i, nil
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
	var i *Item
	var err error
	if i, err = newItem(portUuid, configId, config.Type, config.Name, config.Prams.Copy()); err != nil {
		port.Wallet += config.Cost
		log.Log.Errorf("failed to make item err=[%v]", err)
		return "", fmt.Errorf("failed to make item")
	}
	i.InnerItem.SetPortfolioUuid(portUuid)
	ItemsPortInventory[port.Uuid][i.PortfolioUuid] = i
	Items[i.Uuid] = i

	notification.NewItemNotification(portUuid, i.Type, i.Uuid)
	go port.Update()
	return i.Uuid, nil
}

func (i *Item) DeleteItem() error {
	return DeleteItem(i.Uuid, i.PortfolioUuid, true, true, true)
}

func DeleteItem(uuid, portfolioUuid string, broadcastDelete, lockAcquired, force bool) error {
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

	dbErr := database.Db.DeleteItem(item)
	if dbErr != nil {
		log.Log.Errorf("Failed to delete Item in database uuid=%v err=[%v]", item.Uuid, dbErr)
		if !force {
			return fmt.Errorf("Opps something went wrong 0x01432")
		}
	}

	change.UnregisterChangeDetect(item)
	close(item.UpdateChannel)
	delete(Items, uuid)
	delete(ItemsPortInventory[item.PortfolioUuid], uuid)
	if len(ItemsPortInventory[item.PortfolioUuid]) == 0 {
		delete(ItemsPortInventory, item.PortfolioUuid)
	}
	utils.RemoveUuid(uuid)

	if broadcastDelete {
		sender.SendDeleteObject(portfolioUuid, item)
		wires.ItemsDelete.Offer(item)
	}

	return dbErr
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
 * Transfer an item
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

func UnmarshalJsonItem(itemType, jsonStr string) (InnerItem, error) {
	var item InnerItem
	switch itemType {
	case TradeItemType:
		item = &TradeEffectItem{}
	}

	if err := json.Unmarshal([]byte(jsonStr), &item); err != nil {
		return nil, err
	}
	return item, nil
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
