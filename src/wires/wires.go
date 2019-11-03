package wires

import (
	"github.com/ThisWillGoWell/stock-simulator-server/src/duplicator"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
)

var ItemsNewObjects = duplicator.MakeDuplicator("items-new")
var ItemsUpdate = duplicator.MakeDuplicator("items-update")
var ItemsDelete = duplicator.MakeDuplicator("items-delete")

var PortfolioUpdate = duplicator.MakeDuplicator("portfolio-update")
var PortfolioNewObject = duplicator.MakeDuplicator("new-portfolio")

var UsersNewObject = duplicator.MakeDuplicator("new-users")
var UsersUpdate = duplicator.MakeDuplicator("user-update")

var StocksUpdate = duplicator.MakeDuplicator("valuable-update")
var StocksNewObject = duplicator.MakeDuplicator("new-valuable")

var LedgerUpdate = duplicator.MakeDuplicator("ledger-entries-update")
var LedgerNewObject = duplicator.MakeDuplicator("leger-entries-new")

var NotificationUpdate = duplicator.MakeDuplicator("notification-entries-update")
var NotificationNewObject = duplicator.MakeDuplicator("notification-new")
var NotificationsDelete = duplicator.MakeDuplicator("notification-delete")

var RecordsNewObject = duplicator.MakeDuplicator("records-new")
var BookNewObject = duplicator.MakeDuplicator("book-new")
var BookUpdate = duplicator.MakeDuplicator("book-update")

var EffectsNewObject = duplicator.MakeDuplicator("new-effects")
var EffectsDelete = duplicator.MakeDuplicator("delete-effects")
var EffectsUpdate = duplicator.MakeDuplicator("update-effects")

func ConnectWires() {
	// Enable Copy Mode on all the global new input channels
	UsersUpdate.EnableCopyMode()
	UsersNewObject.EnableCopyMode()
	StocksNewObject.EnableCopyMode()
	PortfolioNewObject.EnableCopyMode()
	LedgerNewObject.EnableCopyMode()
	NotificationNewObject.EnableCopyMode()
	BookNewObject.EnableCopyMode()
	// enable copy mode only user, the rest have copy mode on a channel before
	ItemsUpdate.EnableCopyMode()
	BookUpdate.EnableCopyMode()

	EffectsNewObject.EnableCopyMode()
	EffectsUpdate.EnableCopyMode()

}

func PrintAll() {
	ConnectWires()
	allWires := duplicator.MakeDuplicator("all")
	var out chan interface{}
	out = ItemsNewObjects.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = ItemsUpdate.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = ItemsDelete.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = PortfolioUpdate.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = PortfolioNewObject.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = UsersNewObject.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = UsersUpdate.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = StocksUpdate.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = StocksNewObject.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = LedgerUpdate.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = LedgerNewObject.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = NotificationUpdate.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = NotificationNewObject.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = NotificationsDelete.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = RecordsNewObject.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = BookNewObject.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = BookUpdate.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = EffectsNewObject.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = EffectsDelete.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	out = EffectsUpdate.GetBufferedOutput(100)
	allWires.RegisterInput(out)
	go func() {
		all := allWires.GetBufferedOutput(10000)
		for ele := range all {
			utils.PrintJson(ele)
		}
	}()

}
