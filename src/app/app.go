package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/portfolio"

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
}
type ConfigJson struct {
	Stocks   map[string]JsonStock   `json:"stocks"`
	Accounts map[string]JsonAccount `json:"accounts"`
}

func LoadConfig() {
	configFilePath := os.Getenv("CONFIG_FILE")
	dat, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("error reading config: ", err)
		return
	}
	var config ConfigJson
	err = json.Unmarshal(dat, &config)
	if err != nil {
		fmt.Println("error reading config: ", err)
		return
	}
	for stockId, stockConfig := range config.Stocks {
		_, err = valuable.NewStock(stockId, stockConfig.Name, int64(rand.Intn(10000)), stockConfig.Change.Duration)
		if err != nil {
			fmt.Println("error adding stock: ", err)
		}
	}

	for username, userConfig := range config.Accounts {
		_, err = account.NewUser(username, userConfig.Name, userConfig.Password)
		if err != nil {
			fmt.Println("error adding user: ", err)
		}
	}
	for _, portfolio := range portfolio.Portfolios {
		portfolio.Wallet = 10000000
		for stockId := range valuable.Stocks {
			r := order.MakePurchaseOrder(stockId, portfolio.Uuid, 1)
			<-r.ResponseChannel
			<-time.After(100 * time.Millisecond)
		}
	}

}
