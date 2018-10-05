package items

import (
	"github.com/pkg/errors"
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/notification"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/sender"
)

var ItemTypes = ItemMap()
var ItemsPortInventory = make(map[string]map[string]Item)
var ItemLock = lock.NewLock("item")

const ItemIdentifiableType = "item"

func ItemMap() map[string]ItemType {
	mapp := make(map[string]ItemType)
	mapp[insiderTradingItemType] = InsiderTraderItemType{}
	mapp[mailItemType] = MailItemType{}
	return mapp
}

type ItemType interface {
	GetName() string
	GetType() string
	GetCost() int64
	GetDescription() string
	GetActivateParameters() interface{}
	RequiredLevel() int64
}

type Item interface {
	GetType() string
	GetId() string
	GetItemType() ItemType
	GetUserUuid() string
	GetUuid() string
	SetUserUuid(string)
	Activate(interface{}) (interface{}, error)
	HasBeenUsed() bool
	View() interface{}
	GetUpdateChan() chan interface{}
}

type UserInventory struct {
	UpdateChan chan interface{}
	Items      map[string]string `json:"items" change:"-"`
}

func (*UserInventory) GetType() string {
	return "item_inventory"
}

func makeItem(itemType ItemType, userUuid string) Item {
	switch itemType.(type) {
	case InsiderTraderItemType:
		return newInsiderTradingItem(userUuid)
	case MailItemType:
		return newMailItem(userUuid)
	}
	return nil
}

func BuyItem(portUuid, userUuid, itemName string) (string, error) {

	port, _ := portfolio.GetPortfolio(portUuid)
	itemType, exists := ItemTypes[itemName]
	if !exists {
		return "", errors.New("item type does not exists")
	}
	port.Lock.Acquire("buy item")
	defer port.Lock.Release()
	ItemLock.Acquire("buy-item")
	defer ItemLock.Release()

	if itemType.RequiredLevel() > port.Level {
		return "", errors.New("not high enough level")
	}
	if itemType.GetCost() > port.Wallet {
		return "", errors.New("not enough $$ in wallet")
	}

	port.Wallet -= itemType.GetCost()
	if _, ok := ItemsPortInventory[port.Uuid]; !ok {
		ItemsPortInventory[port.Uuid] = make(map[string]Item)
	}
	newItem := makeItem(itemType, portUuid)
	ItemsPortInventory[port.Uuid][newItem.GetUuid()] = newItem
	change.RegiserPrivateChangeDetect(newItem, newItem.GetUpdateChan())
	sender.GetSender(port.UserUUID).NewObjects.Offer(newItem)
	sender.GetSender(port.UserUUID).Updates.RegisterInput(newItem.GetUpdateChan())

	notification.NewItemNotification(userUuid, itemType.GetName(), newItem.GetId())
	return newItem.GetId(), nil
}

func GetItemsForUser(portfolioUuid string) []Item {
	ItemLock.Acquire("get-Items")
	defer ItemLock.Release()
	items := make([]Item, 0)
	userItems, ok := ItemsPortInventory[portfolioUuid]
	if !ok {
		return items
	}
	for _, item := range userItems {
		items = append(items, item)
	}
	return items
}

func getItem(itemId, portfolioUuid string) (Item, error) {
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

func ViewItem(itemId, userUuid string) (interface{}, error) {
	ItemLock.Acquire("Use Item")
	defer ItemLock.Release()
	item, err := getItem(itemId, userUuid)
	if err != nil {
		return nil, err
	}
	if !item.HasBeenUsed() {
		return nil, errors.New("Item has not been used")
	}
	return item.View(), nil
}

func Use(itemId, portfolioUuid, userUuid string, itemParameters interface{}) (interface{}, error) {
	ItemLock.Acquire("Use Item")
	defer ItemLock.Release()
	item, err := getItem(itemId, portfolioUuid)
	if err != nil {
		return nil, err
	}
	if item.HasBeenUsed() {
		return nil, errors.New("Item has been used")
	}
	val, err := item.Activate(itemParameters)
	if err != nil {
		notification.UsedItemNotification(userUuid, itemId, item.GetItemType().GetName())
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
		ItemsPortInventory[newOwner] = make(map[string]Item)
	}
	ItemsPortInventory[currentOwner][itemId] = item

	delete(ItemsPortInventory[currentOwner], itemId)
	if len(ItemsPortInventory[currentOwner]) == 0 {
		delete(ItemsPortInventory, currentOwner)
	}
	return nil
}
