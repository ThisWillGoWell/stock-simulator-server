package seed

import (
	"encoding/json"
	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"
	"github.com/ThisWillGoWell/stock-simulator-server/src/game/order"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/user"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/valuable"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
	"github.com/ThisWillGoWell/stock-simulator-server/src/web/session"
	"io/ioutil"
	"math/rand"
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
type ObjectsJson struct {
	Stocks   map[string]JsonStock   `json:"stocks"`
	Accounts map[string]JsonAccount `json:"accounts"`
	AutoBuy  bool                   `json:"auto_buy"`
}

func SeedObjects(path string) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Log.Error("error reading file", err)
		return
	}
	var seed ObjectsJson
	if err := json.Unmarshal(data, &seed); err != nil {
		log.Log.Error("error reading seed: ", err)
		return
	}

	stocks := make([]string, 0)
	portfolios := make([]string, 0)

	for stockId, stockConfig := range seed.Stocks {
		stock, err := valuable.NewStock(stockId, stockConfig.Name, int64(rand.Intn(10000)), stockConfig.Change.Duration)
		if err != nil {
			log.Log.Error("error adding stock from seed: ", err)
		} else {
			stocks = append(stocks, stock.Uuid)
		}
	}

	for username, userConfig := range seed.Accounts {
		token, err := user.NewUser(username, userConfig.Name, userConfig.Password)
		if err != nil {
			log.Log.Error("error making user from seed seed: ", err)
		} else {
			uuid, _ := session.GetUserId(token)
			if userConfig.Wallet != 0 {
				portfolio.Portfolios[user.UserList[uuid].PortfolioId].Wallet = userConfig.Wallet
				portfolios = append(portfolios, user.UserList[uuid].PortfolioId)
			}

		}

	}
	if seed.AutoBuy {
		for _, stock := range stocks {
			for _, port := range portfolios {
				order.MakePurchaseOrder(stock, port, 1)
			}
		}
	}

	log.Log.Info("Config done loaded", err)

}
