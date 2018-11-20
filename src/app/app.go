package app

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/stock-simulator-server/src/portfolio"

	"github.com/stock-simulator-server/src/order"

	"github.com/stock-simulator-server/src/log"

	"github.com/stock-simulator-server/src/session"

	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/valuable"

	"github.com/stock-simulator-server/src/account"
)

type JsonStock struct {
	Name   string         `json:"name"`
	Change utils.Duration `json:"change"`
}

type JsonAccount struct {
	Name     string `json:"display_name"`
	Password string `json:"password"`
	Wallet   int64  `json:"wallet"`
}
type ConfigJson struct {
	Stocks   map[string]JsonStock   `json:"stocks"`
	Accounts map[string]JsonAccount `json:"accounts"`
	AutoBuy  bool                   `json:"auto_buy"`
}

func LoadConfig(dat []byte) {
	stocks := make([]string, 0)
	portfolios := make([]string, 0)
	fmt.Println("loading")

	var config ConfigJson
	err := json.Unmarshal(dat, &config)
	if err != nil {
		log.Log.Error("error reading config: ", err)
		return
	}
	for stockId, stockConfig := range config.Stocks {
		stock, err := valuable.NewStock(stockId, stockConfig.Name, int64(rand.Intn(10000)), stockConfig.Change.Duration)
		if err != nil {
			log.Log.Error("error adding stock from config: ", err)
		} else {
			stocks = append(stocks, stock.Uuid)
		}
	}

	for username, userConfig := range config.Accounts {
		token, err := account.NewUser(username, userConfig.Name, userConfig.Password)
		if err != nil {
			log.Log.Error("error making user from config config: ", err)
		} else {
			user, _ := session.GetUserId(token)
			if userConfig.Wallet != 0 {
				portfolio.Portfolios[account.UserList[user].PortfolioId].Wallet = userConfig.Wallet
				portfolios = append(portfolios, account.UserList[user].PortfolioId)
			}

		}

	}
	if config.AutoBuy {
		for _, stock := range stocks {
			for _, port := range portfolios {
				order.MakePurchaseOrder(stock, port, 1)
			}
		}
	}

	log.Log.Info("Config done loaded", err)

}
