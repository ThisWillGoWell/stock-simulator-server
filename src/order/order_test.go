package order

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/account"

	"github.com/ThisWillGoWell/stock-simulator-server/src/valuable"
)

func TestOrder(t *testing.T) {
	Run()
	stock, _ := valuable.NewStock("TEST", "test-stock", 100, time.Second)
	account.NewUser("testuser", "testUser", "pass")
	user, _ := account.GetUser("testuser", "pass")
	po := MakePurchaseOrder(stock.Uuid, user.PortfolioId, 10)
	str, _ := json.Marshal(<-po.ResponseChannel)
	fmt.Println(string(str))
	stock.CurrentPrice = 1000
	po = MakePurchaseOrder(stock.Uuid, user.PortfolioId, -5)
	str, _ = json.Marshal(<-po.ResponseChannel)
	fmt.Println(string(str))
}
