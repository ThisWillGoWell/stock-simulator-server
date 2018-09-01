package app

import (
	"fmt"
	"github.com/stock-simulator-server/src/messages"
	"time"

	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"

	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/client"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/trade"
)

func LoadVars() {
	fmt.Println("running app")
	//make the stocks
	type stockConfig struct {
		id       string
		name     string
		price    float64
		duration time.Duration
	}
	stockConfigs := append(make([]stockConfig, 0),
		stockConfig{"AAPL", "Apple Inc.", 10, time.Minute * 10 },
		stockConfig{"ABBV", "AbbVie Inc.", 10, time.Minute * 10 },
		stockConfig{"ABT", "Abbott Laboratories", 10, time.Minute * 10 },
		stockConfig{"ACN", "Accenture plc", 10, time.Minute * 10 },
		stockConfig{"AGN", "Allergan plc", 10, time.Minute * 10 },
		stockConfig{"AIG", "American International Group Inc.", 10, time.Minute * 10 },
		stockConfig{"ALL", "Allstate Corp.", 10, time.Minute * 10 },
		stockConfig{"AMGN", "Amgen Inc.", 10, time.Minute * 10 },
		stockConfig{"AMZN", "Amazon.com", 10, time.Minute * 10 },
		stockConfig{"AXP", "American Express Inc.", 10, time.Minute * 10 },
		stockConfig{"BA", "Boeing Co.", 10, time.Minute * 10 },
		stockConfig{"BAC", "Bank of America Corp", 10, time.Minute * 10 },
		stockConfig{"BIIB", "Biogen Idec", 10, time.Minute * 10 },
		stockConfig{"BK", "The Bank of New York Mellon", 10, time.Minute * 10 },
		stockConfig{"BKNG", "Booking Holdings", 10, time.Minute * 10 },
		stockConfig{"BLK", "BlackRock Inc", 10, time.Minute * 10 },
		stockConfig{"BMY", "Bristol-Myers Squibb", 10, time.Minute * 10 },
		stockConfig{"BRK.B", "Berkshire Hathaway", 10, time.Minute * 10 },
		stockConfig{"C", "Citigroup Inc", 10, time.Minute * 10 },
		stockConfig{"CAT", "Caterpillar Inc", 10, time.Minute * 10 },
		stockConfig{"CELG", "Celgene Corp", 10, time.Minute * 10 },
		stockConfig{"CHTR", "Charter Communications", 10, time.Minute * 10 },
		stockConfig{"CL", "Colgate-Palmolive Co.", 10, time.Minute * 10 },
		stockConfig{"CMCSA", "Comcast Corporation", 10, time.Minute * 10 },
		stockConfig{"COF", "Capital One Financial Corp.", 10, time.Minute * 10 },
		stockConfig{"COP", "ConocoPhillips", 10, time.Minute * 10 },
		stockConfig{"COST", "Costco", 10, time.Minute * 10 },
		stockConfig{"CSCO", "Cisco Systems", 10, time.Minute * 10 },
		stockConfig{"CVS", "CVS Health", 10, time.Minute * 10 },
		stockConfig{"CVX", "Chevron", 10, time.Minute * 10 },
		stockConfig{"DHR", "Danaher", 10, time.Minute * 10 },
		stockConfig{"DIS", "The Walt Disney Company", 10, time.Minute * 10 },
		stockConfig{"DUK", "Duke Energy", 10, time.Minute * 10 },
		stockConfig{"DWDP", "DowDuPont", 10, time.Minute * 10 },
		stockConfig{"EMR", "Emerson Electric Co.", 10, time.Minute * 10 },
		stockConfig{"EXC", "Exelon", 10, time.Minute * 10 },
		stockConfig{"F", "Ford Motor", 10, time.Minute * 10 },
		stockConfig{"FB", "Facebook", 10, time.Minute * 10 },
		stockConfig{"FDX", "FedEx", 10, time.Minute * 10 },
		stockConfig{"FOX", "21st Century Fox", 10, time.Minute * 10 },
		stockConfig{"FOXA", "21st Century Fox", 10, time.Minute * 10 },
		stockConfig{"GD", "General Dynamics", 10, time.Minute * 10 },
		stockConfig{"GE", "General Electric Co.", 10, time.Minute * 10 },
		stockConfig{"GILD", "Gilead Sciences", 10, time.Minute * 10 },
		stockConfig{"GM", "General Motors", 10, time.Minute * 10 },
		stockConfig{"GOOG", "Alphabet Inc", 10, time.Minute * 10 },
		stockConfig{"GOOGL", "Alphabet Inc", 10, time.Minute * 10 },
		stockConfig{"GS", "Goldman Sachs", 10, time.Minute * 10 },
		stockConfig{"HAL", "Halliburton", 10, time.Minute * 10 },
		stockConfig{"HD", "Home Depot", 10, time.Minute * 10 },
		stockConfig{"HON", "Honeywell", 10, time.Minute * 10 },
		stockConfig{"IBM", "International Business Machines", 10, time.Minute * 10 },
	)

	//Make an exchange //exchanger, _ := exchange.BuildExchange("US")
	//#exchanger.StartExchange()

	//Register stocks with Exchange
	for _, ele := range stockConfigs {
		valuable.NewStock(ele.id, ele.name, ele.price, ele.duration)
	}
	fmt.Println("done adding stocks")

	//start the builder
	//go client.BroadcastMessageBuilder()
	//build and simulate a client
	account.NewUser("Mike", "pass")
	account.NewUser("Will", "pass")
	account.NewUser("Luke", "pass")
	account.NewUser("Chunt", "whip")

	acc, _ := account.GetUser("Will", "pass")
	portfolio.Portfolios[acc.PortfolioId].Wallet = 100000000
	acc2, _ := account.GetUser("Mike", "pass")
	portfolio.Portfolios[acc2.PortfolioId].Wallet = 100000000
	acc3, _ := account.GetUser("Luke", "pass")
	portfolio.Portfolios[acc3.PortfolioId].Wallet = 100000000
	accs := []string{acc.PortfolioId, acc2.PortfolioId, acc3.PortfolioId}
	users := []string{acc.Uuid, acc2.Uuid, acc3.Uuid}

	for id := range valuable.Stocks {
		for i:=0; i<100; i++{
			portId := accs[i%3]
			po2 := order.BuildPurchaseOrder(id, portId, 1)
			trade.Trade(po2)

			client.SendToUser(users[i%3],messages.BuildPurchaseResponse( <-po2.ResponseChannel))
			<-time.After(time.Second * 30)
		}
	}

}
