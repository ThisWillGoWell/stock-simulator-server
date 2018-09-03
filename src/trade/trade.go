package trade

import (
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"
)

var tradeChannel = make(chan *order.PurchaseOrder, 30)

func Trade(o *order.PurchaseOrder) {
	tradeChannel <- o
}

func RunTrader() {
	go func() {
		for purchaseOrder := range tradeChannel {
			trade(purchaseOrder)
		}
	}()
}

// validate and make a trade
// Purchase Order contains a reference to the order
// use failureTrade() and successTrade to send response down channel
// Don't need to a lock around this since the portfolio holds it for that trade
func trade(o *order.PurchaseOrder) {
	//get the stock if it exists
	//valuable.ValuablesLock.EnableDebug()
	//ledger.EntriesLock.EnableDebug()
	valuable.ValuablesLock.Acquire("trade")
	defer valuable.ValuablesLock.Release()

	ledger.EntriesLock.Acquire("trade")
	defer ledger.EntriesLock.Release()

	value, exists := valuable.Stocks[o.ValuableID]
	if !exists {
		order.FailureOrder("asset is not recognized", o)
		return
	}
	//todo possible deadlock
	value.GetLock().Acquire("trade")
	defer value.GetLock().Release()

	port, exists := portfolio.Portfolios[o.PortfolioID]
	if !exists {
		order.FailureOrder("portfolio does not exist, this is very bad", o)
		return
	}

	port.Lock.Acquire("trade")
	defer port.Lock.Release()
	ledgerEntry, ledgerExists := ledger.EntriesStockPortfolio[o.ValuableID][o.PortfolioID]
	if !ledgerExists {
		ledgerEntry = ledger.NewLedgerEntry(o.PortfolioID, o.ValuableID, true)
		port.UpdateInput.RegisterInput(ledgerEntry.UpdateChannel.GetBufferedOutput(10))
	}

	if o.Amount > 0 {
		//we have a buy
		// are there enough shares
		if value.OpenShares < o.Amount {
			order.FailureOrder("not enough open shares", o)
			return
		}
		// does the account have enough money
		costOfTrade := o.Amount * value.GetValue()
		if costOfTrade > port.Wallet {
			order.FailureOrder("not enough $$", o)
			return
		}
		// any other checks?
		// make the trade
		// subtract from open shares
		value.OpenShares -= o.Amount
		//add the holder amount
		ledgerEntry.Amount += o.Amount
		// Update the portfolio with the new ledgerEntry
		port.Wallet -= costOfTrade
		// update the ledger entry to trigger update
		order.SuccessOrder(o)
	} else {
		// we have a sell
		//make sure they have that many shares

		if ledgerEntry.Amount < o.Amount {
			order.FailureOrder("not enough shares", o)
			return
		}
		// make trade
		// add to open shares
		value.OpenShares += o.Amount
		// remove from ledger
		ledgerEntry.Amount -= o.Amount
		costOfTrade := o.Amount * value.GetValue()
		port.Wallet += costOfTrade
		order.SuccessOrder(o)

		}
	if !ledgerExists{
		ledger.NewObjectChannel.Offer(ledgerEntry)
	} else{
		ledgerEntry.UpdateChannel.Offer(ledgerEntry)
	}
	go port.Update()
}
