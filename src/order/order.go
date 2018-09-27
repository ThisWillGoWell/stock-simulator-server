package order

import (
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/valuable"
)

type Order interface {
	rc() chan *Response
}

var orderChannel = make(chan Order, 30)



func Run() {
	go func() {
		for o := range orderChannel {
			switch o.(type) {
			case *PurchaseOrder:
				executeTrade(o.(*PurchaseOrder))
			case *TransferOrder:
				executeTransfer(o.(*TransferOrder))
			}
		}
	}()
}


type TransferOrder struct{
	ReceiverID string `json:"receiver"`
	PortfolioID string 	`json:"giver"`
	Amount int64 `json:"amount"`
	ResponseChannel chan *Response `json:"-"`
}
func (to *TransferOrder) rc()chan *Response{return to.ResponseChannel}


/**
Purchase order represents a single order
*/
type PurchaseOrder struct {
	ValuableID      string         `json:"stock_id"`
	PortfolioID     string         `json:"portfolio"`
	ExchangeID      string         `json:"-"`
	Amount          int64          `json:"amount"`
	ResponseChannel chan *Response `json:"-"`
}
func (po *PurchaseOrder) rc() chan *Response {return po.ResponseChannel}



type Response struct {
	Order   Order 			`json:"order"`
	Success bool           `json:"success"`
	Err     string         `json:"err,omitempty"`
}

// note this does not validate if the stock exists or not, that's done in the trade() function
func MakePurchaseOrder(valuableID, portfolioUUID string, amount int64) *PurchaseOrder {
	po := &PurchaseOrder{
		//ExchangeID:      exchangeID,
		ValuableID:      valuableID,
		PortfolioID:     portfolioUUID,
		Amount:          amount,
		ResponseChannel: make(chan *Response, 1),
	}
	orderChannel <- po
	return po
}

func MakeTransferOrder(giver, reciever string, amount int64) *TransferOrder{
	to := &TransferOrder{
		PortfolioID: giver,
		ReceiverID: reciever,
		Amount: amount,
		ResponseChannel: make(chan *Response, 1),
	}
	orderChannel <- to
	return to
}

func (po *PurchaseOrder) Execute() {

}

func successOrder(o Order) {
	o.rc() <- &Response{
		Order:   o,
		Success: true,
	}
}
func failureOrder(msg string, o Order) {

	o.rc() <- &Response{
		Order:   o,
		Success: false,
		Err:     msg,
	}
}


// validate and make a trade
// Purchase Order contains a reference to the order
// use failureTrade() and successTrade to send response down channel
// Don't need to a lock around this since the portfolio holds it for that trade
func executeTrade(o *PurchaseOrder) {
	//get the stock if it exists
	//valuable.ValuablesLock.EnableDebug()
	//ledger.EntriesLock.EnableDebug()
	valuable.ValuablesLock.Acquire("trade")
	defer valuable.ValuablesLock.Release()

	ledger.EntriesLock.Acquire("trade")
	defer ledger.EntriesLock.Release()

	value, exists := valuable.Stocks[o.ValuableID]
	if !exists {
		failureOrder("asset is not recognized", o)
		return
	}

	value.GetLock().Acquire("trade")
	defer value.GetLock().Release()

	port, exists := portfolio.Portfolios[o.PortfolioID]
	if !exists {
		failureOrder("portfolio does not exist, this is very bad", o)
		return
	}

	port.Lock.Acquire("trade")
	defer port.Lock.Release()
	ledgerEntry, ledgerExists := ledger.EntriesStockPortfolio[o.ValuableID][o.PortfolioID]

	if o.Amount > 0 {
		//we have a buy
		// are there enough shares
		if value.OpenShares < o.Amount {
			failureOrder("not enough open shares", o)
			return
		}
		// does the account have enough money
		costOfTrade := o.Amount * value.GetValue()
		if costOfTrade > port.Wallet {
			failureOrder("not enough $$", o)
			return
		}
		// any other checks?
		// make the trade
		// subtract from open shares
		value.OpenShares -= o.Amount
		// Update the portfolio with the new ledgerEntry
		port.Wallet -= costOfTrade
		// update the ledger entry to trigger update
		if !ledgerExists {
			ledgerEntry = ledger.NewLedgerEntry(o.PortfolioID, o.ValuableID, true)
			port.UpdateInput.RegisterInput(ledgerEntry.UpdateChannel.GetBufferedOutput(10))

		}
		ledgerEntry.InvestmentValue += costOfTrade
		//add the holder amount
		ledgerEntry.Amount += o.Amount
		successOrder(o)
	} else {
		if !ledgerExists {
			failureOrder("not enough shares", o)
			return
		}

		// we have a sell
		//make sure they have that many shares

		amount := o.Amount * -1
		if ledgerEntry.Amount < amount {
			failureOrder("not enough shares", o)
			return
		}

		// make trade
		// add to open shares
		value.OpenShares += amount
		// remove from ledger
		ledgerEntry.Amount -= amount
		costOfTrade := amount * value.GetValue()
		port.Wallet += costOfTrade
		ledgerEntry.InvestmentValue -= costOfTrade
		if ledgerEntry.Amount == 0{
			ledgerEntry.InvestmentValue = 0
		}
		successOrder(o)
	}
	if !ledgerExists{
		ledger.NewObjectChannel.Offer(ledgerEntry)
		port.UpdateInput.RegisterInput(value.UpdateChannel.GetBufferedOutput(10))

	} else{
		ledgerEntry.UpdateChannel.Offer(ledgerEntry)
		// todo remove also
	}
	go value.Update()
	go port.Update()
}

func executeTransfer(o *TransferOrder) {
	portfolio.PortfoliosLock.Acquire("moneyTransfer")
	defer portfolio.PortfoliosLock.Release()
	port, exists := portfolio.Portfolios[o.PortfolioID]
	if !exists{
		failureOrder("giver portfolio not known", o)
	}
	receiver, exists := portfolio.Portfolios[o.ReceiverID]
	if !exists {
		failureOrder("portfolio does not exist, this is very bad", o)
		return
	}
	port.Lock.Acquire("transfer")
	defer port.Lock.Release()
	receiver.Lock.Acquire("transfer")
	defer receiver.Lock.Release()
	if o.Amount <= 0 {
		failureOrder("invalid amount", o)
		return
	}
	if o.Amount > port.Wallet{
		failureOrder("not enough money", o)
		return
	}
	receiver.Wallet += o.Amount
	port.Wallet -= o.Amount
	successOrder(o)
	go port.Update()
	go receiver.Update()
}
