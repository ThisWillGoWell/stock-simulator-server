package app

import (
	"github.com/stock-simulator-server/exchange"
	"time"
	"github.com/stock-simulator-server/valuable"
	"github.com/stock-simulator-server/portfolio"
	"fmt"

	"github.com/stock-simulator-server/account"
	"github.com/stock-simulator-server/client"
	"github.com/stock-simulator-server/messages"
	"encoding/json"
)

func RunApp(){

	//make the stocks
	stock1, _:= valuable.NewStock("CHUNT", "Chunt's Hats", 69,  time.Second * 60)
	stock2, _:= valuable.NewStock("KING", "Paddle King", 10,  time.Second * 5)
	stock3, _:= valuable.NewStock("CBIO", "Sebio's Streaming Services", 10,  time.Second * 30)
	valuable.StartStockStimulation()

	//Make an exchange
	exchanger, _ := exchange.BuildExchange("US")
	exchanger.StartExchange()

	//Register stocks with Exchange
	exchanger.RegisterValuable(stock1, 100)
	exchanger.RegisterValuable(stock2, 100)
	exchanger.RegisterValuable(stock3, 100)

	//Build Some Portfolios
	portfolio.NewPortfolio("1", "Luis Guzman")
	portfolio.NewPortfolio("2", "Big Blacky")

	//start the builder
	go client.BroadcastMessageBuilder()
	//build and simulate a client
	account.NewUser("username", "password")
	go func(){
		return
		rxSim := make(chan string)
		txSim := make(chan string)
		//client.Login("username", "password", txSim, rxSim)
		go func(){
			msg := messages.BaseMessage{
				Action:messages.TradeAction,
				Value:&messages.TradeMessage{
				StockTicker:"CHUNT",
				ExchangeID: "US",
				Amount: 10,
				},
			}
			str, _ := json.Marshal(msg)
			rxSim <- string(str)
		}()
		for msg := range txSim{
			fmt.Println(msg)
		}
	}()

	/*
	po := order.BuildPurchaseOrder("CHUNT", "US", "1", 10)
	exchange.InitiateTrade(po)
	time.Sleep(2 * time.Second)
	po2 := order.BuildPurchaseOrder("KING", "US", "1", 5)
	exchange.InitiateTrade(po2)
	time.Sleep(2 * time.Second)
	po3 := order.BuildPurchaseOrder("CBIO", "US", "1", 1)
	exchange.InitiateTrade(po3)
	time.Sleep(2 * time.Second)
	po4 := order.BuildPurchaseOrder("CBIO", "US", "2", 1)
	exchange.InitiateTrade(po4)
	time.Sleep(2 * time.Second)
	*/
}