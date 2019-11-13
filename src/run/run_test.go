package run

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/money"

	"github.com/ThisWillGoWell/stock-simulator-server/src/order"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/valuable"

	"github.com/ThisWillGoWell/stock-simulator-server/src/record"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/ledger"
)

func TestProspectAll(t *testing.T) {
	App()
	<-time.After(10 * time.Second)
	ValidateRecords(t)
}

func ValidateRecords(t *testing.T) {
	amounts := make(map[string]int64)
	for _, b := range record.GetAllBooks() {
		summ := int64(0)
		for _, rec := range b.ActiveRecords {
			summ += rec.AmountLeft
		}
		amounts[b.LedgerUuid] = summ
	}

	for _, entry := range ledger.GetAllLedgers() {
		if amounts[entry.Uuid] != entry.Amount {
			t.Log(entry.Uuid, entry.Amount, amounts[entry.Uuid], entry.PortfolioId, entry.StockId)
		}
	}
	t.Log("err")
}

func setAllPortfolioWallets(amount int64) {
	for _, portfolio := range portfolio.GetAllPortfolios() {
		portfolio.Wallet = amount
	}
}
func levelUpAllPortfolios(l int64) {
	for _, portfolio := range portfolio.GetAllPortfolios() {
		portfolio.Level = l
	}
}

func setAllAviableStocks(numShares int64) {
	for _, stock := range valuable.GetAllStocks() {
		stock.OpenShares = numShares
	}
}

func TradeStocks(numTrades int) {
	ports := portfolio.GetAllPortfolios()
	stocks := valuable.GetAllStocks()

	for i := 0; i < numTrades; i++ {
		port := ports[i%len(ports)]
		stock := stocks[i%len(stocks)]
		amount := rand.Int63n(200) - 100
		// if the sell amount is greater than the current owned, set the amount = to the amount owned to sell all
		//  make buy if no amount is owned
		if amount < 0 {
			if _, exists := ledger.EntriesPortfolioStock[port.Uuid]; exists {
				entry, exists := ledger.EntriesPortfolioStock[port.Uuid][stock.Uuid]
				if exists {
					if amount*-1 > entry.Amount {
						amount = entry.Amount * -1
					}
				} else {
					amount = amount * -1
				}
			} else {
				amount = amount * -1
			}
		}

		o := order.MakePurchaseOrder(stock.Uuid, port.Uuid, amount)
		<-o.ResponseChannel

	}

}
func TestSimulatorAndValidate(t *testing.T) {
	App()
	<-time.After(10 * time.Second)
	setAllPortfolioWallets(10 * money.Million)
	setAllAviableStocks(10000000)
	TradeStocks(1000000)
	ValidateRecords(t)

}
