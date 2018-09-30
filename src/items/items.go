package items

import (
	"github.com/pkg/errors"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/portfolio"
)

var ItemTypes = ItemMap()
var ItemsUserInventory = make(map[string]map[string]Item)
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

func BuyItem(userUuid, itemName string) error {
	user := account.UserList[userUuid]
	port := portfolio.Portfolios[user.PortfolioId]
	itemType, exists := ItemTypes[itemName]
	if !exists {
		return errors.New("item type does not exists")
	}
	user.Lock.Acquire("buy item")
	defer user.Lock.Release()
	port.Lock.Acquire("buy item")
	defer port.Lock.Release()
	ItemLock.Acquire("buy-item")
	defer ItemLock.Release()

	if itemType.RequiredLevel() > user.Level {
		return errors.New("not high enough level")
	}
	if itemType.GetCost() > port.Wallet {
		return errors.New("not enough $$ in wallet")
	}

	port.Wallet -= itemType.GetCost()
	if _, ok := ItemsUserInventory[user.Uuid]; !ok {
		ItemsUserInventory[user.Uuid] = make(map[string]Item)
	}
	newItem := makeItem(itemType, userUuid)
	ItemsUserInventory[user.Uuid][newItem.GetUuid()] = newItem
	return nil
}

func GetItemsForUser(userUuid string) []Item {
	ItemLock.Acquire("get-Items")
	defer ItemLock.Release()
	items := make([]Item, 0)
	userItems, ok := ItemsUserInventory[userUuid]
	if !ok {
		return items
	}
	for _, item := range userItems {
		items = append(items, item)
	}
	return items
}

func getItem(itemId, userUuid string) (Item, error) {
	userItems, ok := ItemsUserInventory[userUuid]
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

func Use(itemId, userUuid string, itemParameters interface{}) (interface{}, error) {
	ItemLock.Acquire("Use Item")
	defer ItemLock.Release()
	item, err := getItem(itemId, userUuid)
	if err != nil {
		return nil, err
	}
	if item.HasBeenUsed() {
		return nil, errors.New("Item has been used")
	}
	return item.Activate(itemParameters)
}

func TransferItem(currentOwner, newOwner, itemId string) error {
	if _, ok := ItemsUserInventory[currentOwner]; !ok {
		return errors.New("current owner does not own any items")
	}
	item, ok := ItemsUserInventory[currentOwner][itemId]
	if !ok {
		return errors.New("current owner does not have the item id")
	}

	if _, ok := ItemsUserInventory[newOwner]; !ok {
		ItemsUserInventory[newOwner] = make(map[string]Item)
	}
	ItemsUserInventory[currentOwner][itemId] = item

	delete(ItemsUserInventory[currentOwner], itemId)
	if len(ItemsUserInventory[currentOwner]) == 0 {
		delete(ItemsUserInventory, currentOwner)
	}
	return nil
}
