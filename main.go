package main

import (

	"github.com/stock-simulator-server/src/app"
	"github.com/stock-simulator-server/src/web"
	"fmt"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/exchange"
	"github.com/stock-simulator-server/src/client"
	"github.com/stock-simulator-server/src/valuable"
)

func main() {

	go app.RunApp()
	//Wiring of system
	utils.SubscribeUpdateInputs.RegisterInput(portfolio.PortfoliosUpdateChannel.GetOutput())
	utils.SubscribeUpdateInputs.RegisterInput(exchange.ExchangesUpdateChannel.GetOutput())
	utils.SubscribeUpdateInputs.RegisterInput(valuable.ValuableUpdateChannel.GetOutput())

	//this takes the subscribe output and converts it to a message
	client.BroadcastMessageBuilder()
	utils.StartDetectChanges()

	web.StartHandlers()


	fmt.Println("exited!")
}
