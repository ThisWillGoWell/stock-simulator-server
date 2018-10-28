package items

import (
	"testing"
	"time"

	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"
)

func MakeUser(name string) string {
	//user, _ := account.NewUser(name, "test", "pass")
	//port := portfolio.Portfolios[user]
	//port.Wallet = 1000000
	//user.Level = 2
	//return user.Uuid
	return ""

}

func MakeStocks() {
	valuable.NewStock("TEST1", "Test 3 Account", 1000, time.Second*10)
	valuable.NewStock("TEST2", "Test 2 Account", 1000, time.Second*10)
}

func TestInsiderTrading(t *testing.T) {
	//u  := MakeUser("user1")
	MakeStocks()
	//itemId, _ := BuyItem(u, insiderTradingItemType)

	//vals, _ := Use(itemId, u, nil)
	//enc := json.NewEncoder(os.Stdout)
	//enc.Encode(vals)
}

func TestMail(t *testing.T) {
	order.Run()

	u1 := MakeUser("user1")
	u2 := MakeUser("user2")

	//itemId := BuyItem(u1, mailItemType)
	//_, err := Use(itemId, u1, MailItemParameters{To:u2, ShareCount:1000})
	//t.Log(err)
	t.Log(portfolio.Portfolios[account.UserList[u1].PortfolioId].Wallet)
	t.Log(portfolio.Portfolios[account.UserList[u2].PortfolioId].Wallet)
}
