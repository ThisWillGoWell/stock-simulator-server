package items

import (
	"github.com/pkg/errors"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/portfolio"
)

var ItemTypes = make(map[string]ItemType)
var ItemsUserInventory = make(map[string]map[string]Item)
var ItemLock = lock.NewLock("item")
var NewObjectChannel = duplicator.MakeDuplicator("items-new")
var UpdateChannel = duplicator.MakeDuplicator("items-update")

type ItemType interface {
	GetName() string
	GetCost() int64
	GetDescription() string
	GetActivateParameters() interface{}
	RequiredLevel() int64

}

type Item interface {
	GetType() ItemType
	GetUserUuid() string
	GetUuid() string
	Activate() (interface{}, error)
	HasBeenUsed() bool
}

func makeItem(itemType ItemType, userUuid string) Item{
	switch itemType.(type) {
	case InsiderTraderItemType:
		return newInsiderTradingItem(userUuid)
	}
	return nil
}

func BuyItem(userUuid, itemName string) error{
	user := account.UserList[userUuid]
	port := portfolio.Portfolios[user.PortfolioId]
	itemType, exists := ItemTypes[itemName]
	if !exists{
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
		ItemsUserInventory[user.Uuid] =  make(map[string]Item)
	}
	newItem := makeItem(itemType, userUuid)
	ItemsUserInventory[user.Uuid][newItem.GetUuid()] = newItem

	return nil
}

func GetItemsForUser(userUuid string)[]Item{
	ItemLock.Acquire("get-Items")
	defer ItemLock.Release()
	items := make([]Item, 0)
	userItems, ok := ItemsUserInventory[userUuid]
	if !ok{
		return items
	}
	for _, item := range userItems{
		items = append(items, item)
	}
	return items
}


func Use(itemId, userUuid string, itemParameters interface{})(interface{}, error){
	ItemLock.Acquire("Use Item")
	defer ItemLock.Release()
	userItems, ok := ItemsUserInventory[userUuid]
	if !ok{
		return nil, errors.New("user has no items")
	}
	item, ok := userItems[itemId]
	if !ok{
		return nil, errors.New("user does not have that item")
	}
	if item.HasBeenUsed(){
		return nil, errors.New("Item has been used")
	}
	return item.Activate()
}