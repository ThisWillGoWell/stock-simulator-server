package order

import (
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/level"
	"github.com/stock-simulator-server/src/notification"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/record"
	"github.com/stock-simulator-server/src/valuable"
	"github.com/stock-simulator-server/src/wires"
)

type Order interface {
	rc() chan Response
}

type Response interface {
	response()
	IsSuccess() bool
	GetError() string
}

var orderChannel = make(chan Order, 30)

const BaseTaxRate = 0.12
const TradeFee = 0.03

func Run() {
	go func() {
		for o := range orderChannel {
			switch o.(type) {
			case *TradeOrder:
				executeTrade(o.(*TradeOrder))
			case *TransferOrder:
				executeTransfer(o.(*TransferOrder))
			}
		}
	}()
}

type TransferOrder struct {
	ReceiverID      string        `json:"receiver"`
	PortfolioID     string        `json:"giver"`
	Amount          int64         `json:"amount"`
	ResponseChannel chan Response `json:"-"`
}

func (to *TransferOrder) rc() chan Response { return to.ResponseChannel }

/**
Purchase order represents a single order
*/
type TradeOrder struct {
	ValuableID      string        `json:"stock_id"`
	PortfolioID     string        `json:"portfolio"`
	ExchangeID      string        `json:"-"`
	Amount          int64         `json:"amount"`
	ResponseChannel chan Response `json:"-"`
}

func (po *TradeOrder) rc() chan Response { return po.ResponseChannel }

type BasicResponse struct {
	Order        Order   `json:"order"`
	OrderDetails Details `json:"details,omitempty"`
	Success      bool    `json:"success"`
	Err          string  `json:"err,omitempty"`
}

func (*BasicResponse) response()           {}
func (br *BasicResponse) IsSuccess() bool  { return br.Success }
func (br *BasicResponse) GetError() string { return br.Err }

type Details struct {
	SharePrice int64 `json:"share_price"`
	ShareCount int64 `json:"share_count"`
	ShareValue int64 `json:"shares_valuere"`
	Tax        int64 `json:"tax"`
	Fees       int64 `json:"fees"`
	Bonus      int64 `json:"bonus"`
	Result     int64 `json:"result"`
}

// note this does not validate if the stock exists or not, that's done in the trade() function
func MakePurchaseOrder(valuableID, portfolioUUID string, amount int64) *TradeOrder {
	po := &TradeOrder{
		//ExchangeID:      exchangeID,
		ValuableID:      valuableID,
		PortfolioID:     portfolioUUID,
		Amount:          amount,
		ResponseChannel: make(chan Response, 1),
	}
	orderChannel <- po
	return po
}

func MakeTransferOrder(giver, receiver string, amount int64) *TransferOrder {
	to := &TransferOrder{
		PortfolioID:     giver,
		ReceiverID:      receiver,
		Amount:          amount,
		ResponseChannel: make(chan Response, 1),
	}
	orderChannel <- to
	return to
}

func successOrder(o Order, details Details) {
	o.rc() <- &BasicResponse{
		OrderDetails: details,
		Order:        o,
		Success:      true}
}

func failureOrder(msg string, o Order) {
	o.rc() <- &BasicResponse{
		Order:   o,
		Success: false,
		Err:     msg,
	}
}

// validate and make a trade
// Purchase Order contains a reference to the order
// use failureTrade() and successTrade to send response down channel
// Don't need to a lock around this since the portfolio holds it for that trade
func executeTrade(o *TradeOrder) {
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
	details := Details{}
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
		sharesCanOwn := level.Levels[port.Level].MaxSharesStock

		// update the ledger entry to trigger update
		if !ledgerExists {
			if o.Amount > sharesCanOwn {
				failureOrder("can't own that many shares", o)
				return
			}
			ledgerEntry = ledger.NewLedgerEntry(o.PortfolioID, o.ValuableID, true)
			port.UpdateInput.RegisterInput(ledgerEntry.UpdateChannel.GetBufferedOutput(10))
		} else {
			newShareCount := o.Amount + ledgerEntry.Amount
			if newShareCount > sharesCanOwn {
				failureOrder("can't own that many shares", o)
				return
			}
		}
		value.OpenShares -= o.Amount
		// Update the portfolio with the new ledgerEntry
		details = calculateBuyDetails(o, value, port)
		port.Wallet += details.Result
		//add the holder amount
		ledgerEntry.Amount += o.Amount
		successOrder(o, details)
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
		details = calculateSellDetails(o, value, port, ledgerEntry.RecordBookId)
		port.Wallet += details.Result
		successOrder(o, details)
	}
	if !ledgerExists {
		wires.LedgerUpdate.Offer(ledgerEntry)
		port.UpdateInput.RegisterInput(value.UpdateChannel.GetBufferedOutput(1000))
	} else {
		ledgerEntry.UpdateChannel.Offer(ledgerEntry)
	}
	record.NewRecord(ledgerEntry.RecordBookId, details.ShareCount, details.SharePrice, details.Tax, details.Fees, details.Bonus, details.Result)
	notification.DoneTradeNotification(port.UserUUID, value.Uuid, o.Amount)
	go value.Update()
	go port.Update()
}

func CalculateDetails(portfolioUuid, valuableUuid string, amount int64) *BasicResponse {
	order := &TradeOrder{
		ValuableID:  valuableUuid,
		PortfolioID: portfolioUuid,
		Amount:      amount,
	}
	response := &BasicResponse{Order: order}

	v, ok := valuable.Stocks[valuableUuid]
	if !ok {
		response.Err = "valuable id not found"
		return response
	}
	port, ok := portfolio.Portfolios[portfolioUuid]
	if !ok {
		response.Err = "portfolio id not found"
		return response
	}
	v.GetLock().Acquire("calculate-order-details")
	defer v.GetLock().Release()
	port.Lock.Acquire("calculate-order-details")
	defer port.Lock.Release()
	var ledgerEntry *ledger.Entry
	recordUuid := ""

	if order.Amount > 0 {
		response.Success = true
		response.OrderDetails = calculateBuyDetails(order, v, port)
	} else {
		ledgerPortfolio, ledgerExists := ledger.EntriesPortfolioStock[portfolioUuid]
		if ledgerExists {
			ledgerEntry, ledgerExists = ledgerPortfolio[v.Uuid]
		}
		if !ledgerExists {
			response.Err = "can't calculate sell order for ledger that does not exist"
		} else {
			if ledgerEntry.Amount < order.Amount {
				response.Err = "don't own that many shares"

			} else {
				recordUuid = ledgerEntry.RecordBookId
				response.Success = true
				response.OrderDetails = calculateSellDetails(order, v, port, recordUuid)

			}
		}
	}
	return response
}

func calculateBuyDetails(order *TradeOrder, v *valuable.Stock, port *portfolio.Portfolio) Details {
	return Details{
		SharePrice: v.CurrentPrice,
		ShareCount: order.Amount,
		ShareValue: v.CurrentPrice * order.Amount,
		Tax:        0,
		Fees:       0,
		Bonus:      0,
		Result:     v.CurrentPrice * order.Amount * -1,
	}
}

func calculateSellDetails(order *TradeOrder, v *valuable.Stock, port *portfolio.Portfolio, recordUuid string) Details {
	d := Details{
		SharePrice: v.CurrentPrice,
		ShareCount: order.Amount,
		ShareValue: v.CurrentPrice * order.Amount * -1,
	}
	principle := record.GetPrinciple(recordUuid, order.Amount*-1)
	pbt := d.ShareValue - principle
	taxes := 0.0
	if pbt > 0 {
		taxes = float64(pbt) * BaseTaxRate
	}
	fees := int64(float64(d.ShareValue) * TradeFee)
	d.Tax = int64(taxes)
	d.Fees = int64(fees)
	d.Result = d.ShareValue - d.Tax - d.Fees + d.Bonus
	return d
}

func executeTransfer(o *TransferOrder) {
	if o.ReceiverID == o.PortfolioID {
		failureOrder("cant transfer to and from same account", o)
		return
	}
	port, exists := portfolio.Portfolios[o.PortfolioID]
	if !exists {
		failureOrder("giver portfolio not known", o)
		return
	}

	if port.Level == 0 {
		failureOrder("need to be level 1 to transfer money", o)
		return
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

	if o.Amount > port.Wallet {
		failureOrder("not enough money", o)
		return
	}
	receiver.Wallet += o.Amount
	port.Wallet -= o.Amount
	successOrder(o, Details{})
	notification.SendMoneyTradeNotification(port.UserUUID, receiver.UserUUID, o.Amount)
	go port.Update()
	go receiver.Update()
}
