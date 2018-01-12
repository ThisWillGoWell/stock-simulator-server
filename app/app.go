package app

import (
	"stock-server/exchange"
	"stock-server/utils"
	"time"
	"stock-server/valuable"
	"stock-server/portfolio"
	"fmt"
	"stock-server/order"
)


var UpdateChannel = utils.MakeDuplicator()

func RunApp(){
	//make the stocks
	stock1, _:= valuable.NewStock("CHUNT", "Chunt's Hats", 69,  time.Second * 10)
	stock2, _:= valuable.NewStock("KING", "Paddle King", 10,  time.Second * 2)
	stock3, _:= valuable.NewStock("CBIO", "Sebio's Streaming Services", 10,  time.Second)

	//Make an exchange
	exchanger, _ := exchange.BuildExchange("US")
	exchanger.StartExchange()

	//Register stocks with Exchange
	exchanger.RegisterValuable(stock1, 100)
	exchanger.RegisterValuable(stock2, 100)
	exchanger.RegisterValuable(stock3, 100)

	//Build Some Portfolios
	port, _ := portfolio.NewPortfolio("1")

	//make a purchase order
	po := order.BuildPurchaseOrder(valuable.Valuables["CHUNT"], port, 10)
	exchange.Exchanges["US"].InitiateTrade(po)

	result := <- po.ResponseChannel
	fmt.Println(result)

}