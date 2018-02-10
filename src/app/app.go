package app

import (
	"fmt"
	"github.com/stock-simulator-server/src/exchange"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"
	"time"

	"encoding/json"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/messages"
	"github.com/stock-simulator-server/src/order"
)

func RunApp() {
	fmt.Println("running app")
	//make the stocks
	type stockConfig struct {
		id       string
		name     string
		price    float64
		duration time.Duration
	}

	stockConfigs := append(make([]stockConfig, 0),
		stockConfig{"CHUNT", "Chunt's Hats", 69, time.Second * 45},
		stockConfig{"KING", "Paddle King", 10, time.Second * 30},
		stockConfig{"CBIO", "Sebio's Streaming Services", 10, time.Minute * 1},
		stockConfig{"OW", "Overwatch", 10, time.Minute * 2},
		stockConfig{"SCOTT", "Michael Scott Paper Company ", 10, time.Minute * 3},
		stockConfig{"DM", "Dunder Milf ", 10, time.Minute * 4},
		stockConfig{"GWEN", "", 10, time.Minute * 5},
		stockConfig{"CHU", "Chu Supply", 10, time.Minute * 4},
		stockConfig{"SWEET", "Sweet Sweet Tea", 10, time.Minute * 3},
		stockConfig{"TRAP", "❤ Trap 4 Life", 10, time.Minute * 2},
		stockConfig{"FIG", "Figgis Agency", 10, time.Minute * 2},
		stockConfig{"ZONE", "Danger Zone", 10, time.Minute * 1},
		stockConfig{"PLNX", "Planet Express", 10, time.Minute * 2},
		stockConfig{"MOM", "Mom's Friendly Robot Company", 10, time.Minute * 3},
	)
	valuable.StartStockStimulation()

	//Make an exchange
	exchanger, _ := exchange.BuildExchange("US")
	exchanger.StartExchange()

	//Register stocks with Exchange
	for _, ele := range stockConfigs {
		stock, _ := valuable.NewStock(ele.id, ele.name, ele.price, ele.duration)
		exchanger.RegisterValuable(stock, 100)
	}

	//Build Some Portfolios
	portfolio.NewPortfolio("1", "Luis Guzman")
	portfolio.NewPortfolio("2", "Big Blacky")

	//start the builder
	//go client.BroadcastMessageBuilder()
	//build and simulate a client
	account.NewUser("username", "password")
	go func() {
		return
		rxSim := make(chan string)
		txSim := make(chan string)
		//client.Login("username", "password", txSim, rxSim)
		go func() {
			msg := messages.BaseMessage{
				Action: messages.TradeAction,
				Msg: &messages.TradeMessage{
					StockTicker: "CHUNT",
					ExchangeID:  "US",
					Amount:      10,
				},
			}
			str, _ := json.Marshal(msg)
			rxSim <- string(str)
		}()
		for msg := range txSim {
			fmt.Println(msg)
		}
	}()

	po2 := order.BuildPurchaseOrder("KING", "US", "1", 5)
	exchange.InitiateTrade(po2)
	time.Sleep(10 * time.Second)
	po3 := order.BuildPurchaseOrder("CBIO", "US", "1", 1)
	exchange.InitiateTrade(po3)
	time.Sleep(10 * time.Second)
	po4 := order.BuildPurchaseOrder("CBIO", "US", "2", 1)
	exchange.InitiateTrade(po4)
	time.Sleep(10 * time.Second)

}