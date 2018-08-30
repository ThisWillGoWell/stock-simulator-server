package wires

import (
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/client"
	"github.com/stock-simulator-server/src/database"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"
)

func ConnectWires(diableDb bool) {
	// Enable Copy Mode on all the global new input channels
	account.NewObjectChannel.EnableCopyMode()
	valuable.NewObjectChannel.EnableCopyMode()
	portfolio.NewObjectChannel.EnableCopyMode()
	ledger.NewObjectChannel.EnableCopyMode()

	// enable copy mode only account, the rest have copy mode on a channel before
	account.UpdateChannel.EnableCopyMode()

	//Build a new object channel
	newObjectChannels := duplicator.MakeDuplicator("new-objects")

	newObjectChannels.RegisterInput(account.NewObjectChannel.GetOutput())
	newObjectChannels.RegisterInput(valuable.NewObjectChannel.GetOutput())
	newObjectChannels.RegisterInput(portfolio.NewObjectChannel.GetOutput())
	newObjectChannels.RegisterInput(ledger.NewObjectChannel.GetOutput())

	client.Updates.RegisterInput(newObjectChannels.GetOutput())
	if !diableDb {
		database.DatabseWriter.RegisterInput(newObjectChannels.GetOutput())
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
