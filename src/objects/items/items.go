package items

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/database"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"

	"github.com/ThisWillGoWell/stock-simulator-server/src/merge"

	"github.com/ThisWillGoWell/stock-simulator-server/src/change"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/log"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/notification"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/sender"
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
	models.Item
	UpdateChannel chan interface{} `json:"-"`
}

type i2 struct {
	Uuid          string    `json:"uuid"`
	Name          string    `json:"name"`
	ConfigId      string    `json:"config"`
	Type          string    `json:"type"`
	PortfolioUuid string    `json:"portfolio_uuid"`
	CreateTime    time.Time `json:"create_time"`
	InnerItem     InnerItem `json:"-" change:"inner"`
}

func (i *Item) GetId() string {
	return i.Uuid
}

func (*Item) GetType() string {
	return ItemIdentifiableType
}

func newItem(portfolioUuid, configId, itemType, name string, innerItem interface{}) (i *Item, err error) {
	return MakeItem(models.Item{
		Name:          name,
		ConfigId:      configId,
		Uuid:          id.SerialUuid(),
		PortfolioUuid: portfolioUuid,
		Type:          itemType,
		CreateTime:    time.Now(),
		InnerItem:     innerItem,
	})
}

func MakeItem(i models.Item) (*Item, error) {
	switch i.InnerItem.(type) {
	case string:
		var err error
		if  i.InnerItem, err = UnmarshalJsonItem(i.Type,  i.InnerItem.(string)); err != nil {
			return nil, err
		}
	}
	item := &Item{
		Item: i,
		UpdateChannel: make(chan interface{}),
	}
	if err := sender.RegisterChangeUpdate(item.PortfolioUuid, item.UpdateChannel); err != nil {
		return nil, err
	}

	if _, ok := ItemsPortInventory[item.PortfolioUuid]; !ok {
		ItemsPortInventory[item.PortfolioUuid] = make(map[string]*Item)
	}
	item.InnerItem.(InnerItem).SetParentItemUuid(item.Uuid)

	if err := change.RegisterPrivateChangeDetect(item, item.UpdateChannel); err != nil {
		return nil, err
	}

	id.RegisterUuid(item.Uuid, i)
	ItemsPortInventory[item.PortfolioUuid][item.Uuid] = i
	Items[item.Uuid] = item

	return item, nil
}

func BuyItem(portUuid, configId string) (string, error) {
	portfolio.PortfoliosLock.Acquire("buy-item")
	port, ok := portfolio.Portfolios[portUuid]
	if !ok {
		portfolio.PortfoliosLock.Release()
		return "", fmt.Errorf("portfolio not found")
	}
	port.Lock.Acquire("buy-item")
	defer port.Lock.Release()
	portfolio.PortfoliosLock.Release()

	config, found := validItems[configId]
	if !found {
		return "", errors.New("config not found")
	}

	if config.RequiredLevel > port.Level {
		return "", errors.New("not high enough level")
	}
	if config.Cost > port.Wallet {
		return "", errors.New("not enough $$ in wallet")
	}

	ItemLock.Acquire("buy-item")
	defer ItemLock.Release()

	if _, ok := ItemsPortInventory[port.Uuid]; !ok {
		ItemsPortInventory[port.Uuid] = make(map[string]*Item)
	}

	var i *Item
	var err error
	if i, err = newItem(portUuid, configId, config.Type, config.Name, config.Prams.Copy()); err != nil {
		log.Log.Errorf("failed to make item err=[%v]", err)
		return "", fmt.Errorf("failed to make item")
	}

	i.InnerItem.(InnerItem).SetPortfolioUuid(portUuid)
	ItemsPortInventory[port.Uuid][i.Uuid] = i
	Items[i.Uuid] = i

	port.Wallet -= config.Cost

	notification.NotificationLock.Acquire("buy item")
	defer notification.NotificationLock.Release()

	note := notification.NewItemNotification(portUuid, i.Type, i.Uuid)

	//commit them all to the database
	if dbErr := database.Db.Execute([]interface{}{note.Notification, i.Item, port.Portfolio}, nil); dbErr != nil {
		// undo the item buy
		delete(ItemsPortInventory[port.Uuid], i.Uuid)
		if len(ItemsPortInventory[port.Uuid]) == 0 {
			delete(ItemsPortInventory, port.Uuid)
		}
		notification.DeleteNotification(note.Uuid, true)
		port.Wallet += config.Cost
		deleteItem(i)
		log.Log.Errorf("failed to buy item err=[%v]", err)
		return "", fmt.Errorf("oops! something went wrong")
	}
	return i.Uuid, nil
}

func (i *Item) DeleteItem() error {
	return database.Db.Execute(nil, []interface{}{i})
}

func DeleteItem(uuid string) error {

	ItemLock.Acquire("delete-item")
	defer ItemLock.Release()
	// delete from the database
	item, exists := Items[uuid]
	if !exists {
		log.Log.Errorf("got a delete for a item that does not exists")
		return nil
	}
	if err := database.Db.Execute(nil, []interface{}{item}); err != nil {
		log.Log.Errorf("failed to delete item err=[%v]", err)
		return fmt.Errorf("opps! something went wrong")
	}
	// we can delete the item now
	deleteItem(item)
	sender.SendDeleteObject(item.PortfolioUuid, item)
	return nil
}

func deleteItem(item *Item) {
	change.UnregisterChangeDetect(item)
	close(item.UpdateChannel)
	delete(Items, item.Uuid)
	delete(ItemsPortInventory[item.PortfolioUuid], item.Uuid)
	if len(ItemsPortInventory[item.PortfolioUuid]) == 0 {
		delete(ItemsPortInventory, item.PortfolioUuid)
	}
	id.RemoveUuid(item.Uuid)
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

	val, err := item.InnerItem.(InnerItem).Activate(itemParameters)
	if err != nil {
		if err := notification.UsedItemNotification(portfolioUuid, itemId, item.Type); err != nil {
			log.Log.Errorf("failed to make %s used-item notification for %s err=[%v]", portfolioUuid, err)
		}
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
