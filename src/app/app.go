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
)

func LoadVars() {
	fmt.Println("running app")
	//make the stocks
	type stockConfig struct {
		id       string
		name     string
		price    int64
		duration time.Duration
	}
	stockConfigs := append(make([]stockConfig, 0),
		stockConfig{"CHUNT", "Slut Hut", 1000, time.Second },
		stockConfig{"ABBV", "AbbVie Inc.", 1000, time.Second * 1 },
		stockConfig{"ABT", "Abbott Laboratories", 1000, time.Second * 1 },
		stockConfig{"ACN", "Accenture plc", 1000, time.Second * 1 },
		stockConfig{"AGN", "Allergan plc", 1000, time.Second * 1 },
		stockConfig{"AIG", "American International Group Inc.", 1000, time.Second * 1 },
		stockConfig{"ALL", "Allstate Corp.", 1000, time.Second * 1 },
		stockConfig{"AMGN", "Amgen Inc.", 1000, time.Second * 1 },
		stockConfig{"AMZN", "Amazon.com", 1000, time.Second * 1 },
		stockConfig{"AXP", "American Express Inc.", 1000, time.Second * 1 },
		stockConfig{"BA", "Boeing Co.", 1000, time.Second * 1 },
		stockConfig{"BAC", "Bank of America Corp", 1000, time.Second * 1 },
		stockConfig{"BIIB", "Biogen Idec", 1000, time.Second * 1 },
		stockConfig{"BK", "The Bank of New York Mellon", 1000, time.Second * 1 },
		stockConfig{"BKNG", "Booking Holdings", 1000, time.Second * 1 },
		stockConfig{"BLK", "BlackRock Inc", 1000, time.Second * 1 },
		stockConfig{"BMY", "Bristol-Myers Squibb", 1000, time.Second * 1 },
		stockConfig{"BRK.B", "Berkshire Hathaway", 1000, time.Second * 1 },
		stockConfig{"C", "Citigroup Inc", 1000, time.Second * 1 },
		stockConfig{"CAT", "Caterpillar Inc", 1000, time.Second * 1 },
		stockConfig{"CELG", "Celgene Corp", 1000, time.Second * 1 },
		stockConfig{"CHTR", "Charter Communications", 1000, time.Second * 1 },
		stockConfig{"CL", "Colgate-Palmolive Co.", 1000, time.Second * 1 },
		stockConfig{"CMCSA", "Comcast Corporation", 1000, time.Second * 1 },
		stockConfig{"COF", "Capital One Financial Corp.", 1000, time.Second * 1 },
		stockConfig{"COP", "ConocoPhillips", 1000, time.Second * 1 },
		stockConfig{"COST", "Costco", 1000, time.Second * 1 },
		stockConfig{"CSCO", "Cisco Systems", 1000, time.Second * 1 },
		stockConfig{"CVS", "CVS Health", 1000, time.Second * 1 },
		stockConfig{"CVX", "Chevron", 1000, time.Second * 1 },
		stockConfig{"DHR", "Danaher", 1000, time.Second * 1 },
		stockConfig{"DIS", "The Walt Disney Company", 1000, time.Second * 1 },
		stockConfig{"DUK", "Duke Energy", 1000, time.Second * 1 },
		stockConfig{"DWDP", "DowDuPont", 1000, time.Second * 1 },
		stockConfig{"EMR", "Emerson Electric Co.", 1000, time.Second * 1 },
		stockConfig{"EXC", "Exelon", 1000, time.Second * 1 },
		stockConfig{"F", "Ford Motor", 1000, time.Second * 1 },
		stockConfig{"FB", "Facebook", 1000, time.Second * 1 },
		stockConfig{"FDX", "FedEx", 1000, time.Second * 1 },
		stockConfig{"FOX", "21st Century Fox", 1000, time.Second * 1 },
		stockConfig{"FOXA", "21st Century Fox", 1000, time.Second * 1 },
		stockConfig{"GD", "General Dynamics", 1000, time.Second * 1 },
		stockConfig{"GE", "General Electric Co.", 1000, time.Second * 1 },
		stockConfig{"GILD", "Gilead Sciences", 1000, time.Second * 1 },
		stockConfig{"GM", "General Motors", 1000, time.Second * 1 },
		stockConfig{"GOOG", "Alphabet Inc", 1000, time.Second * 1 },
		stockConfig{"GOOGL", "Alphabet Inc", 1000, time.Second * 1 },
		stockConfig{"GS", "Goldman Sachs", 1000, time.Second * 1 },
		stockConfig{"HAL", "Halliburton", 1000, time.Second * 1 },
		stockConfig{"HD", "Home Depot", 1000, time.Second * 1 },
		stockConfig{"HON", "Honeywell", 1000, time.Second * 1 },
		stockConfig{"IBM", "International Business Machines", 1000, time.Second * 1 },
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
			po2 := order.MakePurchaseOrder(id, portId, 1)
			client.SendToUser(users[i%3],messages.BuildPurchaseResponse( <-po2.ResponseChannel))
			<-time.After(time.Second * 30)
			to := order.MakeTransferOrder(portId, accs[(i+1)%3], 10)
			client.SendToUser(users[i%3],messages.BuildTransferResponse( <-to.ResponseChannel))

		}
	}

}
