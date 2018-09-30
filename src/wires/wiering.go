package wires

import (
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/database"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/valuable"
)

var ItemsNewChannel = duplicator.MakeDuplicator("items-new")
var ItemsUpdateChannel = duplicator.MakeDuplicator("items-update")

var PortfolioUpdateChannel = duplicator.MakeDuplicator("portfolio-update")
var PortfolioNewObjectChannel = duplicator.MakeDuplicator("new-portfolio")

var UsersNewObjectChannel = duplicator.MakeDuplicator("new-users")
var UsersUpdateChannel = duplicator.MakeDuplicator("user-update")


var GlobalNewObjects = duplicator.MakeDuplicator("global-new-objects")
var GlobalDeletes = duplicator.MakeDuplicator("global-deletes")
var GlobalNotifications = duplicator.MakeDuplicator("global-notifications")
var GlobalUpdates = duplicator.MakeDuplicator("global-new-objects")
var Globals = duplicator.MakeDuplicator("global-broadcast")

// change detector
var PublicSubscribeInputs = duplicator.MakeDuplicator("subscribe-inputs")
var PublicSubscribeChagneOutputs

func ConnectWires(diableDb bool) {
	// Enable Copy Mode on all the global new input channels
	UsersNewObjectChannel.EnableCopyMode()
	valuable.NewObjectChannel.EnableCopyMode()
	PortfolioNewObjectChannel.EnableCopyMode()
	ledger.NewObjectChannel.EnableCopyMode()

	// enable copy mode only account, the rest have copy mode on a channel before
	UsersUpdateChannel.EnableCopyMode()

	GlobalUpdates.RegisterInput(ItemsUpdateChannel.GetOutput())
	GlobalUpdates.RegisterInput(valuable.NewObjectChannel.GetOutput())
	GlobalUpdates.RegisterInput(PortfolioNewObjectChannel.GetOutput())
	GlobalUpdates.RegisterInput(ledger.NewObjectChannel.GetOutput())

	if !diableDb {
		database.DatabseWriter.RegisterInput(GlobalUpdates.GetOutput())
	}

	// register changes to change detector
	change.SubscribeUpdateInputs.RegisterInput(portfolio.UpdateChannel.GetOutput())
	change.SubscribeUpdateInputs.RegisterInput(ledger.UpdateChannel.GetOutput())
	change.SubscribeUpdateInputs.RegisterInput(valuable.UpdateChannel.GetOutput())
	change.SubscribeUpdateInputs.RegisterInput(account.UpdateChannel.GetOutput())

	//register output
	client.Updates.RegisterInput(change.SubscribeUpdateOutput.GetOutput())
	if !diableDb {
		database.DatabseWriter.RegisterInput(change.SubscribeUpdateOutput.GetOutput())
	}

}
