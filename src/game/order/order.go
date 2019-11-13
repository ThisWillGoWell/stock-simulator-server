package order

import (
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/sender"

	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"
	"github.com/ThisWillGoWell/stock-simulator-server/src/database"
	"github.com/ThisWillGoWell/stock-simulator-server/src/game/level"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/effect"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/ledger"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/notification"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/record"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/valuable"
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

func Run() {
	go func() {
		for o := range orderChannel {
			switch o.(type) {
			case *TradeOrder:
				executeTrade(o.(*TradeOrder))
			case *TransferOrder:
				executeTransfer(o.(*TransferOrder))
			case *ProspectOrder:
				calculateDetails(o.(*ProspectOrder))
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

type ProspectOrder struct {
	ValuableID      string        `json:"stock_id"`
	PortfolioID     string        `json:"portfolio"`
	ExchangeID      string        `json:"-"`
	Amount          int64         `json:"amount"`
	ResponseChannel chan Response `json:"-"`
}

func (po *ProspectOrder) rc() chan Response { return po.ResponseChannel }

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
	SharePrice    int64    `json:"share_price"`
	ShareCount    int64    `json:"share_count"`
	ShareValue    int64    `json:"shares_valuere"`
	Tax           int64    `json:"tax"`
	Fees          int64    `json:"fees"`
	Bonus         int64    `json:"bonus"`
	Result        int64    `json:"result"`
	ActiveEffects []string `json:"active_effects"`
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

func MakeProspect(valuableID, portfolioUUID string, amount int64) *ProspectOrder {
	po := &ProspectOrder{
		//ExchangeID:      exchangeID,
		ValuableID:      valuableID,
		PortfolioID:     portfolioUUID,
		Amount:          amount,
		ResponseChannel: make(chan Response, 1),
	}
	orderChannel <- po
	return po
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
	var err error
	valuable.ValuablesLock.Acquire("trade")
	value, exists := valuable.Stocks[o.ValuableID]
	if !exists {
		valuable.ValuablesLock.Release()
		failureOrder("asset is not recognized", o)
		return
	}
	// lock the object for the rest of the trade
	value.GetLock().Acquire("trade")
	defer value.GetLock().Release()
	valuable.ValuablesLock.Release()

	portfolio.PortfoliosLock.Acquire("trade")
	port, exists := portfolio.Portfolios[o.PortfolioID]
	if !exists {
		portfolio.PortfoliosLock.Release()
		failureOrder("portfolio does not exist, this is very bad", o)
		return
	}
	//lock the portfolio for the rest of the trade
	port.Lock.Acquire("trade")
	portfolio.PortfoliosLock.Release()

	// get the ledger or make a new one

	ledger.EntriesLock.Acquire("trade")
	ledger.NewEntriesLock.Acquire("trade")
	ledgerEntry, ledgerExists := ledger.EntriesStockPortfolio[o.ValuableID][o.PortfolioID]
	if !ledgerExists {
		defer ledger.NewEntriesLock.Release()
	} else {
		ledger.NewEntriesLock.Release()
		ledgerEntry.Lock.Acquire("trade")
		defer ledgerEntry.Lock.Release()
	}
	ledger.EntriesLock.Release()
	// no need to aquire lock here, nothing changed or added
	tradeEffects, activeEffects := effect.TotalTradeEffect(port.Uuid)

	details := Details{}
	if o.Amount > 0 {
		//we have a buy
		// are there enough shares
		if value.OpenShares < o.Amount {
			failureOrder("not enough open shares", o)
			return
		}
		details = calculateBuyDetails(o.Amount, value, tradeEffects, activeEffects)

		// does the user have enough money
		if details.Result > port.Wallet {
			failureOrder("not enough $$", o)
			return
		}
		totalShareCountOwned := o.Amount
		if ledgerExists {
			totalShareCountOwned += ledgerEntry.Amount
		}
		if totalShareCountOwned > level.Levels[port.Level].MaxSharesStock {
			failureOrder("can't own that many shares", o)
			return
		}

	} else {
		if !ledgerExists {
			failureOrder("not enough shares", o)
			return
		}
		// we have a sell
		//make sure they have that many shares
		if ledgerEntry.Amount >= o.Amount {
			failureOrder("not enough shares", o)
			return
		}
		details = calculateSellDetails(o.Amount, value, ledgerEntry.RecordBookId, tradeEffects, activeEffects)
	}

	// update/create the ledger
	if !ledgerExists {
		// todo acquire lock when its made
		ledgerEntry, err = ledger.NewLedgerEntry(o.PortfolioID, o.ValuableID)
		if err != nil {
			log.Log.Errorf("err during ledger create err=[%v]", err)
			failureOrder(fmt.Sprintf("Opps! Something went wrong 0x043"), o)
			return
		}
	} else {
		ledgerEntry.Amount += o.Amount
	}

	// make the record
	r, book := record.NewRecord(ledgerEntry.RecordBookId, details.ShareCount, details.SharePrice, details.Tax, details.Fees, details.Bonus, details.Result)
	// make the notification
	note := notification.DoneTradeNotification(port.Uuid, value.Uuid, o.Amount)

	// make the trade
	port.Wallet += details.Result
	value.OpenShares -= o.Amount

	if dbErr := database.Db.Execute([]interface{}{port.Portfolio, value.Stock, ledgerEntry.Ledger, r.Record, note}, nil); dbErr != nil {
		// undo the trade
		port.Wallet -= details.Result
		value.OpenShares += o.Amount
		ledgerEntry.Amount -= o.Amount
		record.DeleteRecord(r.Uuid, true)
		if ledgerEntry.Amount == 0 {
			ledger.DeleteLedger(ledgerEntry, true)
		}
		notification.DeleteNotification(note.Uuid, true)
		record.DeleteRecord(r.Uuid, true)
		log.Log.Errorf("failed to make trade err=[%v]", err)
		failureOrder("Oops! something went wrong!", o)
		return
	}

	// we have committed the stuff to the database, offer to all downstream listeners
	if !ledgerExists {
		port.UpdateInput.RegisterInput(value.UpdateChannel.GetBufferedOutput(100))
		wires.LedgerNewObject.Offer(ledgerEntry)
		wires.BookNewObject.Offer(book)
	} else {
		ledgerEntry.UpdateChannel.Offer(ledgerEntry)
	}
	go port.Update()
	go value.Update()
	// send the new objects
	sender.SendNewObject(port.Uuid, note)
	wires.RecordsNewObject.Offer(r)

	successOrder(o, details)
	go port.Update()
}

func calculateDetails(order *ProspectOrder) {
	response := &BasicResponse{Order: order}

	valuable.ValuablesLock.Acquire("prospect")
	v, ok := valuable.Stocks[order.ValuableID]
	if !ok {
		valuable.ValuablesLock.Release()
		response.Err = "valuable id not found"
		return
	}
	// lock for rest
	v.GetLock().Acquire("prospect")
	defer v.GetLock().Release()
	valuable.ValuablesLock.Release()

	portfolio.PortfoliosLock.Acquire("prospect")
	port, ok := portfolio.Portfolios[order.PortfolioID]
	if !ok {
		portfolio.PortfoliosLock.Release()
		response.Err = "portfolio id not found"
		return
	}
	port.Lock.Acquire("calculate-order-details")
	defer port.Lock.Release()
	portfolio.PortfoliosLock.Release()

	tradeEffect, activeEffects := effect.TotalTradeEffect(order.PortfolioID)

	var ledgerEntry *ledger.Entry
	recordUuid := ""

	if order.Amount > 0 {
		response.Success = true
		response.OrderDetails = calculateBuyDetails(order.Amount, v, tradeEffect, activeEffects)
		order.ResponseChannel <- response
		return
	}
	ledgerPortfolio, ledgerExists := ledger.EntriesPortfolioStock[order.PortfolioID]
	if ledgerExists {
		ledgerEntry, ledgerExists = ledgerPortfolio[v.Uuid]
	}
	if !ledgerExists || ledgerEntry == nil {
		response.Err = "can't calculate sell order for ledger that does not exist"
		order.ResponseChannel <- response
		return
	}
	if ledgerEntry.Amount < order.Amount*-1 {
		response.Err = "don't own that many shares"
		order.ResponseChannel <- response
		return
	}
	recordUuid = ledgerEntry.RecordBookId
	response.Success = true
	response.OrderDetails = calculateSellDetails(order.Amount, v, recordUuid, tradeEffect, activeEffects)
	order.ResponseChannel <- response

}

func calculateBuyDetails(amount int64, v *valuable.Stock, tradeEffect *effect.TradeEffect, activeTradeEffects []string) Details {

	d := Details{
		SharePrice:    v.CurrentPrice,
		ShareCount:    amount,
		ShareValue:    v.CurrentPrice * amount,
		Tax:           0,
		Fees:          int64(float64(*tradeEffect.BuyFeeAmount) * *tradeEffect.BuyFeeMultiplier),
		Bonus:         0,
		Result:        v.CurrentPrice * amount * -1,
		ActiveEffects: activeTradeEffects,
	}
	d.Result = d.ShareValue*-1 - d.Fees
	return d
}

func calculateSellDetails(amount int64, v *valuable.Stock, recordUuid string, tradeEffect *effect.TradeEffect, activeTradeEffects []string) Details {
	d := Details{
		SharePrice:    v.CurrentPrice,
		ShareCount:    amount,
		ShareValue:    v.CurrentPrice * amount * -1,
		ActiveEffects: activeTradeEffects,
	}
	principle := record.GetPrinciple(recordUuid, amount*-1)
	pbt := d.ShareValue - principle
	taxes := 0.0
	if pbt > 0 {
		taxes = float64(pbt) * *tradeEffect.TaxPercent
		d.Bonus = int64(float64(pbt) * *tradeEffect.BonusProfitMultiplier)
	}
	fees := int64(float64(*tradeEffect.SellFeeAmount) * *tradeEffect.SellFeeMultiplier)
	d.Tax = int64(taxes)
	d.Fees = int64(fees)
	d.Result = d.ShareValue - d.Tax - d.Fees + d.Bonus

	return d
}

func executeTransfer(o *TransferOrder) {
	if o.ReceiverID == o.PortfolioID {
		failureOrder("cant transfer to and from same user", o)
		return
	}
	portfolio.PortfoliosLock.Acquire("transfer money")
	port, exists := portfolio.Portfolios[o.PortfolioID]
	if !exists {
		portfolio.PortfoliosLock.Release()
		failureOrder("giver portfolio not known", o)
		return
	}

	receiver, exists := portfolio.Portfolios[o.ReceiverID]
	if !exists {
		portfolio.PortfoliosLock.Release()
		failureOrder("receiver portolio not found und", o)
		return
	}
	// lock the lock both portfolios for the rest of the trade
	port.Lock.Acquire("transfer")
	defer port.Lock.Release()
	receiver.Lock.Acquire("transfer")
	defer receiver.Lock.Release()
	portfolio.PortfoliosLock.Release()

	if port.Level == 0 {
		failureOrder("need to be level 1 to transfer money", o)
		return
	}

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

	notification.NotificationLock.Acquire("transfer")
	defer notification.NotificationLock.Release()

	n1, n2 := notification.SendMoneyTradeNotification(port.Uuid, receiver.Uuid, o.Amount)

	// commit the changes to the database
	if dbErr := database.Db.Execute([]interface{}{n1.Notification, n2.Notification, receiver.Portfolio, port.Portfolio}, nil ); dbErr != nil {
		// undo the transfer
		receiver.Wallet -= o.Amount
		port.Wallet += o.Amount
		notification.DeleteNotification(n1.Uuid, true)
		notification.DeleteNotification(n2.Uuid, true)
		log.Log.Errorf("failed to write to database send money err=[%v]", dbErr)
		failureOrder("Oops! something went wrong", o)
		return
	}

	successOrder(o, Details{})
	sender.SendNewObject(port.Uuid, n1)
	sender.SendNewObject(receiver.Uuid, n2)
	go port.Update()
	go receiver.Update()
}
