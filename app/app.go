package app

import (
	"github.com/stock-simulator-server/exchange"
	"github.com/stock-simulator-server/utils"
	"time"
	"github.com/stock-simulator-server/valuable"
	"github.com/stock-simulator-server/portfolio"
	"fmt"
	"github.com/stock-simulator-server/order"
	"encoding/json"
	"github.com/stock-simulator-server/messages"
)


var UpdateChannel = utils.MakeDuplicator()

const portfolioUpdateName  = "portfolio"
const ledgerUpdateName  = "ledger"
const stockUpdateName = "stock"

func RunApp(){
	txMessage := utils.MakeDuplicator()

	portfolioMessageBuilder := utils.MakeDuplicator()
	stockMessageBuilder := utils.MakeDuplicator()
	ledgerMessageBuilder := utils.MakeDuplicator()

	//something that take the portfolio update and turns it into a portfolio Update message
	go func(){
		updateChannel := portfolioMessageBuilder.GetOutput()
		for update := range updateChannel{
			txMessage.Offer(messages.NewUpdateMessage(portfolioUpdateName, update))
			}
	}()
	//something that take the stock update and turns it into a portfolio Update message
	go func(){
		updateChannel := stockMessageBuilder.GetOutput()
		for update := range updateChannel{
			txMessage.Offer(messages.NewUpdateMessage(stockUpdateName, update))
		}
	}()

	//something that take the portfolio update and turns it into a portfolio Update message
	go func(){
		updateChannel := ledgerMessageBuilder.GetOutput()
		for update := range updateChannel{
			txMessage.Offer(messages.NewUpdateMessage(ledgerUpdateName, update))
		}
	}()

	// this is simulating a tx websocket line
	go func(){
		output := txMessage.GetOutput()
		for update := range output{
			jsonPrintStr, err := json.Marshal(update)
			if err != nil{
				fmt.Println(err)
			}
			fmt.Println(string(jsonPrintStr))
		}
	}()

	//make the stocks
	stock1, _:= valuable.NewStock("CHUNT", "Chunt's Hats", 69,  time.Second * 60)
	stock2, _:= valuable.NewStock("KING", "Paddle King", 10,  time.Second * 5)
	stock3, _:= valuable.NewStock("CBIO", "Sebio's Streaming Services", 10,  time.Second * 30)
	valuable.StartStockStimulation()
	//register the stocks update to the stock message builder
	stockMessageBuilder.RegisterInput(stock1.GetUpdateChannel().GetOutput())
	stockMessageBuilder.RegisterInput(stock2.GetUpdateChannel().GetOutput())
	stockMessageBuilder.RegisterInput(stock3.GetUpdateChannel().GetOutput())

	//Make an exchange
	exchanger, _ := exchange.BuildExchange("US")
	exchanger.StartExchange()
	// register the exchange's ledger update channel to the ledger message builder
	ledgerMessageBuilder.RegisterInput(exchanger.LedgerUpdateChannel.GetOutput())
	//Register stocks with Exchange
	exchanger.RegisterValuable(stock1, 100)
	exchanger.RegisterValuable(stock2, 100)
	exchanger.RegisterValuable(stock3, 100)

	//Build Some Portfolios
	port, _ := portfolio.NewPortfolio("1", "Luis Guzman")
	port2, _ := portfolio.NewPortfolio("2", "Big Blacky")
	portfolioMessageBuilder.RegisterInput(port.UpdateChannel.GetOutput())
	portfolioMessageBuilder.RegisterInput(port2.UpdateChannel.GetOutput())
	//make a purchase order
	po := order.BuildPurchaseOrder(valuable.Valuables["CHUNT"], port, 10)
	exchange.Exchanges["US"].InitiateTrade(po)
	time.Sleep(2 * time.Second)
	po2 := order.BuildPurchaseOrder(valuable.Valuables["KING"], port, 5)
	exchange.Exchanges["US"].InitiateTrade(po2)
	time.Sleep(2 * time.Second)
	po3 := order.BuildPurchaseOrder(valuable.Valuables["CBIO"], port, 1)
	exchange.Exchanges["US"].InitiateTrade(po3)
	time.Sleep(2 * time.Second)
	po4 := order.BuildPurchaseOrder(valuable.Valuables["CBIO"], port2, 1)
	exchange.Exchanges["US"].InitiateTrade(po4)
	time.Sleep(2 * time.Second)

	result := <- po.ResponseChannel
	fmt.Println(result)

}