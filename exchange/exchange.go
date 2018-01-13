package exchange

import (
	"github.com/stock-simulator-server/utils"
	"github.com/stock-simulator-server/valuable"
	"errors"
	"github.com/stock-simulator-server/order"
	"github.com/stock-simulator-server/portfolio"
	"fmt"
)

var Exchanges = make(map[string]*Exchange)
var ExchangesLock = utils.NewLock("exchanges")

const TOLERANCE = 0.000001

type Exchange struct {
	name string
	ledger map[string ]*ledgerEntry
	tradeChannel *utils.ChannelDuplicator
	lock *utils.Lock
	LedgerUpdateChannel * utils.ChannelDuplicator
}

type ledgerEntry struct {
	ValuableID		string `json:"id"`
	ExchangeName string `json:"exchange"`
	Holders    map[string]float64 `json:"holders"`
	OpenShares float64            `json:"open_shares"`
}

func BuildExchange(name string) (*Exchange, error){
	ExchangesLock.Acquire("build-exchange")
	defer ExchangesLock.Release()

	if _ , exists := Exchanges[name]; exists {
		return nil, errors.New("exchange already exists")
	}

	exchange := &Exchange{
		name: name,
		ledger:               	make(map[string]*ledgerEntry),
		tradeChannel:         	utils.MakeDuplicator(),
		lock: 					utils.NewLock(fmt.Sprintf("exchange-%s", name)),
		LedgerUpdateChannel: 	utils.MakeDuplicator(),
	}
	Exchanges[name] = exchange
	return exchange, nil
}
// #########################
//			Stocks
// #########################

func  (exchange *Exchange) RegisterValuable(valuable valuable.Valuable, amount float64) error {
	exchange.lock.Acquire("register-valuable")
	defer exchange.lock.Release()
	if _, exists := exchange.ledger[valuable.GetID()]; exists {
		return errors.New("valuable already exists in exchange")
	}
	exchange.ledger[valuable.GetID()] = &ledgerEntry{
		ValuableID: valuable.GetID(),
		ExchangeName: exchange.name,
		OpenShares: amount,
		Holders:    make( map[string]float64),
	}
	return nil
}


// #########################
//			Portfolio
// #########################

func (exchange *Exchange)GetHoldingsCount(portfolio *portfolio.Portfolio, valuable valuable.Valuable) float64{
	amount, exists := exchange.ledger[valuable.GetID()].Holders[portfolio.UUID]
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
	o.Valuable.GetLock().Acquire("trade")
	defer o.Valuable.GetLock().Release()
	o.Portfolio.Lock.Acquire("trade")
	defer o.Portfolio.Lock.Release()

	info, ok := exchange.ledger[o.Valuable.GetID()]
	if !ok {
		order.FailureOrder("stock is not known", o)
		return
	}
	if o.Amount > 0{
		//we have a buy
		// are there enough shares
		if info.OpenShares < o.Amount{
			order.FailureOrder("not enough shares", o)
			return
		}
		// does the user have enough money
		costOfTrade := o.Amount * o.Valuable.GetValue()
		if costOfTrade > o.Portfolio.Wallet {
			order.FailureOrder("not enough $$", o)
			return
		}
		// any other checks?

		// make the trade
		// subtract from open shares
		info.OpenShares -= o.Amount
		//add the holder amount
		_, exists := info.Holders[o.Portfolio.UUID]
		if ! exists {
			info.Holders[o.Portfolio.UUID] = o.Amount
		}else{
			info.Holders[o.Portfolio.UUID] += o.Amount
		}
		// Update the portfolio with the new info
		o.Portfolio.TradeUpdate(o.Valuable, info.Holders[o.Portfolio.UUID], costOfTrade)
		order.SuccessOrder(o)
	} else {
		// we have a sell
		//make sure they have that many shares
		currentOwns, exists := info.Holders[o.Portfolio.UUID]
		if ! exists{
			order.FailureOrder("does not own any of stock", o)
			return
		}
		if currentOwns < o.Amount{
			order.FailureOrder("not enough shares", o)
			return
		}
		// make trade
		moneyReturned := o.Amount + o.Valuable.GetValue()
		// add to open shares
		info.OpenShares += o.Amount
		// remove from ledger
		info.Holders[o.Portfolio.UUID] -= o.Amount
		if info.Holders[o.Portfolio.UUID] < TOLERANCE {
			delete(info.Holders, o.Portfolio.UUID)
			o.Portfolio.TradeUpdate(o.Valuable,0, moneyReturned)
		}else{
			o.Portfolio.TradeUpdate(o.Valuable, info.Holders[o.Portfolio.UUID], moneyReturned)
		}
		order.SuccessOrder(o)
	}
	exchange.LedgerUpdateChannel.Offer(info)

}