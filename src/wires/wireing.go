package wires

import (
	"github.com/stock-simulator-server/src/duplicator"
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

func ConnectWires() {
	// Enable Copy Mode on all the global new input channels
	UsersNewObject.EnableCopyMode()
	StocksNewObject.EnableCopyMode()
	PortfolioNewObject.EnableCopyMode()
	LedgerNewObject.EnableCopyMode()
	NotificationNewObject.EnableCopyMode()
	BookNewObject.EnableCopyMode()
	// enable copy mode only account, the rest have copy mode on a channel before
	UsersUpdate.EnableCopyMode()
	ItemsUpdate.EnableCopyMode()
	BookUpdate.EnableCopyMode()

}
