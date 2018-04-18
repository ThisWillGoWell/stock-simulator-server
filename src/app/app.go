package app

import (
	"fmt"
	"time"

	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"

	"encoding/json"

	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/client"
	"github.com/stock-simulator-server/src/messages"
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
		stockConfig{"AAPL", "Apple Inc.", 10, time.Minute},
		stockConfig{"ABBV", "AbbVie Inc.", 10, time.Minute},
		stockConfig{"ABT", "Abbott Laboratories", 10, time.Minute},
		stockConfig{"ACN", "Accenture plc", 10, time.Minute},
		stockConfig{"AGN", "Allergan plc", 10, time.Minute},
		stockConfig{"AIG", "American International Group Inc.", 10, time.Minute},
		stockConfig{"ALL", "Allstate Corp.", 10, time.Minute},
		stockConfig{"AMGN", "Amgen Inc.", 10, time.Minute},
		stockConfig{"AMZN", "Amazon.com", 10, time.Minute},
		stockConfig{"AXP", "American Express Inc.", 10, time.Minute},
		stockConfig{"BA", "Boeing Co.", 10, time.Minute},
		stockConfig{"BAC", "Bank of America Corp", 10, time.Minute},
		stockConfig{"BIIB", "Biogen Idec", 10, time.Minute},
		stockConfig{"BK", "The Bank of New York Mellon", 10, time.Minute},
		stockConfig{"BKNG", "Booking Holdings", 10, time.Minute},
		stockConfig{"BLK", "BlackRock Inc", 10, time.Minute},
		stockConfig{"BMY", "Bristol-Myers Squibb", 10, time.Minute},
		stockConfig{"BRK.B", "Berkshire Hathaway", 10, time.Minute},
		stockConfig{"C", "Citigroup Inc", 10, time.Minute},
		stockConfig{"CAT", "Caterpillar Inc", 10, time.Minute},
		stockConfig{"CELG", "Celgene Corp", 10, time.Minute},
		stockConfig{"CHTR", "Charter Communications", 10, time.Minute},
		stockConfig{"CL", "Colgate-Palmolive Co.", 10, time.Minute},
		stockConfig{"CMCSA", "Comcast Corporation", 10, time.Minute},
		stockConfig{"COF", "Capital One Financial Corp.", 10, time.Minute},
		stockConfig{"COP", "ConocoPhillips", 10, time.Minute},
		stockConfig{"COST", "Costco", 10, time.Minute},
		stockConfig{"CSCO", "Cisco Systems", 10, time.Minute},
		stockConfig{"CVS", "CVS Health", 10, time.Minute},
		stockConfig{"CVX", "Chevron", 10, time.Minute},
		stockConfig{"DHR", "Danaher", 10, time.Minute},
		stockConfig{"DIS", "The Walt Disney Company", 10, time.Minute},
		stockConfig{"DUK", "Duke Energy", 10, time.Minute},
		stockConfig{"DWDP", "DowDuPont", 10, time.Minute},
		stockConfig{"EMR", "Emerson Electric Co.", 10, time.Minute},
		stockConfig{"EXC", "Exelon", 10, time.Minute},
		stockConfig{"F", "Ford Motor", 10, time.Minute},
		stockConfig{"FB", "Facebook", 10, time.Minute},
		stockConfig{"FDX", "FedEx", 10, time.Minute},
		stockConfig{"FOX", "21st Century Fox", 10, time.Minute},
		stockConfig{"FOXA", "21st Century Fox", 10, time.Minute},
		stockConfig{"GD", "General Dynamics", 10, time.Minute},
		stockConfig{"GE", "General Electric Co.", 10, time.Minute},
		stockConfig{"GILD", "Gilead Sciences", 10, time.Minute},
		stockConfig{"GM", "General Motors", 10, time.Minute},
		stockConfig{"GOOG", "Alphabet Inc", 10, time.Minute},
		stockConfig{"GOOGL", "Alphabet Inc", 10, time.Minute},
		stockConfig{"GS", "Goldman Sachs", 10, time.Minute},
		stockConfig{"HAL", "Halliburton", 10, time.Minute},
		stockConfig{"HD", "Home Depot", 10, time.Minute},
		stockConfig{"HON", "Honeywell", 10, time.Minute},
		stockConfig{"IBM", "International Business Machines", 10, time.Minute},
		stockConfig{"INTC", "Intel Corporation", 10, time.Minute},
		stockConfig{"JNJ", "Johnson & Johnson Inc", 10, time.Minute},
		stockConfig{"JPM", "JP Morgan Chase & Co", 10, time.Minute},
		stockConfig{"KHC", "Kraft Heinz", 10, time.Minute},
		stockConfig{"KMI", "Kinder Morgan Inc/DE", 10, time.Minute},
		stockConfig{"KO", "The Coca-Cola Company", 10, time.Minute},
		stockConfig{"LLY", "Eli Lilly and Company", 10, time.Minute},
		stockConfig{"LMT", "Lockheed-Martin", 10, time.Minute},
		stockConfig{"LOW", "Lowe's", 10, time.Minute},
		stockConfig{"MA", "MasterCard Inc", 10, time.Minute},
		stockConfig{"MCD", "McDonald's Corp", 10, time.Minute},
		stockConfig{"MDLZ", "Mondelēz International", 10, time.Minute},
		stockConfig{"MDT", "Medtronic Inc.", 10, time.Minute},
		stockConfig{"MET", "Metlife Inc.", 10, time.Minute},
		stockConfig{"MMM", "3M Company", 10, time.Minute},
		stockConfig{"MO", "Altria Group", 10, time.Minute},
		stockConfig{"MON", "Monsanto", 10, time.Minute},
		stockConfig{"MRK", "Merck & Co.", 10, time.Minute},
		stockConfig{"MS", "Morgan Stanley", 10, time.Minute},
		stockConfig{"MSFT", "Microsoft", 10, time.Minute},
		stockConfig{"NEE", "NextEra Energy", 10, time.Minute},
		stockConfig{"NKE", "Nike", 10, time.Minute},
		stockConfig{"ORCL", "Oracle Corporation", 10, time.Minute},
		stockConfig{"OXY", "Occidental Petroleum Corp.", 10, time.Minute},
		stockConfig{"PEP", "Pepsico Inc.", 10, time.Minute},
		stockConfig{"PFE", "Pfizer Inc", 10, time.Minute},
		stockConfig{"PG", "Procter & Gamble Co", 10, time.Minute},
		stockConfig{"PM", "Phillip Morris International", 10, time.Minute},
		stockConfig{"PYPL", "PayPal Holdings", 10, time.Minute},
		stockConfig{"QCOM", "Qualcomm Inc.", 10, time.Minute},
		stockConfig{"RTN", "Raytheon Company", 10, time.Minute},
		stockConfig{"SBUX", "Starbucks Corporation", 10, time.Minute},
		stockConfig{"SLB", "Schlumberger", 10, time.Minute},
		stockConfig{"SO", "Southern Company", 10, time.Minute},
		stockConfig{"SPG", "Simon Property Group, Inc.", 10, time.Minute},
		stockConfig{"T", "AT&T Inc", 10, time.Minute},
		stockConfig{"TGT", "Target Corp.", 10, time.Minute},
		stockConfig{"TWX", "Time Warner Inc.", 10, time.Minute},
		stockConfig{"TXN", "Texas Instruments", 10, time.Minute},
		stockConfig{"UNH", "UnitedHealth Group Inc.", 10, time.Minute},
		stockConfig{"UNP", "Union Pacific Corp.", 10, time.Minute},
		stockConfig{"UPS", "United Parcel Service Inc", 10, time.Minute},
		stockConfig{"USB", "US Bancorp", 10, time.Minute},
		stockConfig{"UTX", "United Technologies Corp", 10, time.Minute},
		stockConfig{"V", "Visa Inc.", 10, time.Minute},
		stockConfig{"VZ", "Verizon Communications Inc", 10, time.Minute},
		stockConfig{"WBA", "Walgreens Boots Alliance", 10, time.Minute},
		stockConfig{"WFC", "Wells Fargo", 10, time.Minute},
		stockConfig{"WMT", "Wal-Mart", 10, time.Minute},
		stockConfig{"XOM", "Exxon Mobil Corp", 10, time.Minute},
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

	go func() {
		return
		rxSim := make(chan string)
		txSim := make(chan string)
		//client.Login("username", "password", txSim, rxSim)
		go func() {
			msg := messages.BaseMessage{
				Action: messages.TradeAction,
				Msg: &messages.TradeMessage{
					StockId:    "CHUNT",
					ExchangeID: "US",
					Amount:     10,
				},
			}
			str, _ := json.Marshal(msg)
			rxSim <- string(str)
		}()
		for msg := range txSim {
			fmt.Println(msg)
		}
	}()

	acc, _ := account.GetUser("Will", "pass")
	portfolio.Portfolios[acc.PortfolioId].Wallet = 1000000

	for id, _ := range valuable.Stocks {
		po2 := order.BuildPurchaseOrder(id, acc.PortfolioId, 5)
		trade.Trade(po2)
		client.BroadcastMessages.Offer(<-po2.ResponseChannel)
		<-time.After(time.Second * 10)
	}

}
