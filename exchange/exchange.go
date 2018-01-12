package exchange

import (
	"stock-server/utils"

	"stock-server/valuable"
	"errors"
	"stock-server/order"
	"stock-server/portfolio"
	"fmt"
)

var Exchanges = make(map[string]*Exchange)
var ExchangesLock = utils.NewLock()

const TOLERANCE = 0.000001

type Exchange struct {
	name string
	ledger map[valuable.Valuable ]*ledgerEntry
	tradeChannel *utils.ChannelDuplicator
	portfolios map[*portfolio.Portfolio]bool
	lock *utils.Lock
}

type ledgerEntry struct {
	valuable   valuable.Valuable
	holders    map[*portfolio.Portfolio]float64
	openShares float64
}

func BuildExchange(name string) (*Exchange, error){
	ExchangesLock.Acquire()
	defer ExchangesLock.Release()

	if _ , exists := Exchanges[name]; exists {
		return nil, errors.New("exchange already exists")
	}

	exchange := &Exchange{
		name: name,
		ledger:               make(map[valuable.Valuable]*ledgerEntry),
		tradeChannel:         utils.MakeDuplicator(),
		lock: utils.NewLock(),
	}
	Exchanges[name] = exchange
	return exchange, nil
}
// #########################
//			Stocks
// #########################

func  (exchange *Exchange) RegisterValuable(valuable valuable.Valuable, amount float64) error {
	exchange.lock.Acquire()
	defer exchange.lock.Release()
	if _, exists := exchange.ledger[valuable]; exists {
		return errors.New("valuable already exists in exchange")
	}
	exchange.ledger[valuable] = buildLedgerEntry(valuable, amount)
	return nil
}

func buildLedgerEntry(valuable valuable.Valuable, totalShares float64) *ledgerEntry{
	return &ledgerEntry{
		valuable: valuable,
		openShares: totalShares,
		holders: make( map[*portfolio.Portfolio]float64),
	}
}

// #########################
//			Portfolio
// #########################

func (exchange *Exchange)GetHoldingsCount(portfolio *portfolio.Portfolio, valuable valuable.Valuable) float64{
	amount, exists := exchange.ledger[valuable].holders[portfolio]
	if !exists{
		return 0
	}
	return amount
}


// #########################
// 			Trade
// #########################




func (exchange *Exchange) InitiateTrade(o *order.PurchaseOrder){
	fmt.Println(o.Amount)
	fmt.Println(exchange.lock)
	exchange.tradeChannel.Offer(o)
}

func (exchange *Exchange) StartExchange() {
	incomingTrades := exchange.tradeChannel.GetOutput()
	go func(){
		for purchaseOrder := range incomingTrades {
			exchange.trade(purchaseOrder.(*order.PurchaseOrder))
		}
	}()

}


// validate and make a trade
// Purchase Order contains a reference to the order
// use failureTrade() and successTrade to send response down channel
// Don't need to a lock around this since the portfolio holds it for that trade
func (exchange *Exchange) trade(o *order.PurchaseOrder) {
	o.Valuable.GetLock().Acquire()
	defer o.Valuable.GetLock().Release()
	o.Portfolio.Lock.Acquire()
	defer o.Portfolio.Lock.Release()

	info, ok := exchange.ledger[o.Valuable]
	if !ok {
		order.FailureOrder("stock is not known", o)
		return
	}
	if o.Amount > 0{
		//we have a buy
		// are there enough shares
		if info.openShares < o.Amount{
			order.FailureOrder("not enough shares", o)
			return
		}
		// does the user have enough money
		costOfTrade := o.Amount * info.valuable.GetValue()
		if costOfTrade > o.Portfolio.Wallet {
			order.FailureOrder("not enough $$", o)
			return
		}
		// any other checks?

		// make the trade
		// subtract from open shares
		info.openShares -= o.Amount
		//add the holder amount
		_, exists := info.holders[o.Portfolio]
		if ! exists {
			info.holders[o.Portfolio] = o.Amount
		}else{
			info.holders[o.Portfolio] += o.Amount
		}
		// subtract from the wallet
		o.Portfolio.UpdateWallet(-1 * costOfTrade)
		order.SuccessOrder(o)
	} else {
		// we have a sell
		//make sure they have that many shares
		currentOwns, exists := info.holders[o.Portfolio]
		if ! exists{
			order.FailureOrder("does not own any of stock", o)
		}
		if currentOwns < o.Amount{
			order.FailureOrder("not enough shares", o)
		}
		// make trade
		// add the cash into the wallet
		o.Portfolio.Wallet += o.Amount + info.valuable.GetValue()
		// add to open shares
		info.openShares += o.Amount
		// remove from ledger
		info.holders[o.Portfolio] -= o.Amount
		if info.holders[o.Portfolio] < TOLERANCE {
			delete(info.holders, o.Portfolio)
			o.Portfolio.UpdateLedger(o.Valuable,0)
		}else{
			o.Portfolio.UpdateLedger(o.Valuable, info.holders[o.Portfolio])
		}
		order.SuccessOrder(o)
	}

}