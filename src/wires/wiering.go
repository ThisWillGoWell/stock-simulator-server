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

var GlobalNewObjects = duplicator.MakeDuplicator("global-new-objects")
var GlobalDeletes = duplicator.MakeDuplicator("global-deletes")
var GlobalNotifications = duplicator.MakeDuplicator("global-notifications")
var GlobalUpdates = duplicator.MakeDuplicator("global-new-objects")
var Globals = duplicator.MakeDuplicator("global-broadcast")

// change detector
var PublicSubscribeInputs = duplicator.MakeDuplicator("subscribe-inputs")
var PublicSubscribeChagneOutputs = duplicator.MakeDuplicator("public-subscribe-updates")

func ConnectWires(diableDb bool) {
	// Enable Copy Mode on all the global new input channels
	UsersNewObject.EnableCopyMode()
	StocksNewObject.EnableCopyMode()
	PortfolioNewObject.EnableCopyMode()
	LedgerNewObject.EnableCopyMode()

	// enable copy mode only account, the rest have copy mode on a channel before
	UsersUpdate.EnableCopyMode()

	GlobalUpdates.RegisterInput(ItemsUpdate.GetOutput())
	GlobalUpdates.RegisterInput(StocksUpdate.GetOutput())
	GlobalUpdates.RegisterInput(PortfolioUpdate.GetOutput())
	GlobalUpdates.RegisterInput(LedgerUpdate.GetOutput())

}
