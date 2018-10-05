package wires

import (
	"github.com/stock-simulator-server/src/duplicator"
)

var ItemsNewObjects = duplicator.MakeDuplicator("items-new")
var ItemsUpdate = duplicator.MakeDuplicator("items-update")

var PortfolioUpdate = duplicator.MakeDuplicator("portfolio-update")
var PortfolioNewObject = duplicator.MakeDuplicator("new-portfolio")

var UsersNewObject = duplicator.MakeDuplicator("new-users")
var UsersUpdate = duplicator.MakeDuplicator("user-update")

var StocksUpdate = duplicator.MakeDuplicator("valuable-update")
var StocksNewObject = duplicator.MakeDuplicator("new-valuable")

var LedgerUpdate = duplicator.MakeDuplicator("ledger-entries-update")
var LedgerNewObject = duplicator.MakeDuplicator("leger-entries-new")

var NotificationUpdate = duplicator.MakeDuplicator("notification-entries-update")
var NotificationNewObject = duplicator.MakeDuplicator("notification-entries-new")

var GlobalNewObjects = duplicator.MakeDuplicator("global-new-objects")
var GlobalDeletes = duplicator.MakeDuplicator("global-deletes")
var GlobalNotifications = duplicator.MakeDuplicator("global-notifications")
var GlobalUpdates = duplicator.MakeDuplicator("global-updates")
var Globals = duplicator.MakeDuplicator("global-broadcast")

func ConnectWires() {
	// Enable Copy Mode on all the global new input channels
	UsersNewObject.EnableCopyMode()
	StocksNewObject.EnableCopyMode()
	PortfolioNewObject.EnableCopyMode()
	LedgerNewObject.EnableCopyMode()
	NotificationNewObject.EnableCopyMode()

	// enable copy mode only account, the rest have copy mode on a channel before
	UsersUpdate.EnableCopyMode()

	GlobalNewObjects.RegisterInput(UsersNewObject.GetBufferedOutput(10000))
	GlobalNewObjects.RegisterInput(StocksNewObject.GetBufferedOutput(10000))
	GlobalNewObjects.RegisterInput(PortfolioNewObject.GetBufferedOutput(10000))
	GlobalNewObjects.RegisterInput(LedgerNewObject.GetBufferedOutput(10000))

	GlobalUpdates.RegisterInput(ItemsUpdate.GetBufferedOutput(10000))
	GlobalUpdates.RegisterInput(StocksUpdate.GetBufferedOutput(10000))
	GlobalUpdates.RegisterInput(PortfolioUpdate.GetBufferedOutput(10000))
	GlobalUpdates.RegisterInput(LedgerUpdate.GetBufferedOutput(10000))
	GlobalUpdates.RegisterInput(UsersUpdate.GetBufferedOutput(10000))
	GlobalUpdates.RegisterInput(NotificationUpdate.GetBufferedOutput(10000))

}
