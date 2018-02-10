package exchange

import (
	"errors"
	"fmt"
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/valuable"
)

const ObjectType = "exchange_ledger"

var Exchanges = make(map[string]*Exchange)
var ExchangesLock = utils.NewLock("exchanges")
var ExchangesUpdateChannel = utils.MakeDuplicator()

const TOLERANCE = 0.000001

type Exchange struct {
	name                string
	Ledger              map[string]*ledgerEntry
	tradeChannel        *utils.ChannelDuplicator
	lock                *utils.Lock
	LedgerUpdateChannel *utils.ChannelDuplicator
}

type ledgerEntry struct {
	ExchangeName string             `json:"exchange"`
	Name         string             `json:"name"`
	Holders      map[string]float64 `json:"holders" change:"-"`
	OpenShares   float64            `json:"open_shares" change:"-"`
}

func (ledger *ledgerEntry) GetId() string {
	return ledger.ExchangeName + "/" + ledger.Name
}

func (ledger *ledgerEntry) GetType() string {
	return ObjectType
}

func BuildExchange(name string) (*Exchange, error) {
	ExchangesLock.Acquire("build-exchange")
	defer ExchangesLock.Release()

	if _, exists := Exchanges[name]; exists {
		return nil, errors.New("exchange already exists")
	}

	exchange := &Exchange{
		name:                name,
		Ledger:              make(map[string]*ledgerEntry),
		tradeChannel:        utils.MakeDuplicator(),
		lock:                utils.NewLock(fmt.Sprintf("exchange-%s", name)),
		LedgerUpdateChannel: utils.MakeDuplicator(),
	}
	Exchanges[name] = exchange
	ExchangesUpdateChannel.RegisterInput(exchange.LedgerUpdateChannel.GetOutput())
	return exchange, nil
}

// #########################
//			Stocks
// #########################

func (exchange *Exchange) RegisterValuable(valuable valuable.Valuable, amount float64) error {
	exchange.lock.Acquire("register-valuable")
	defer exchange.lock.Release()
	if _, exists := exchange.Ledger[valuable.GetId()]; exists {
		return errors.New("valuable already exists in exchange")
	}
	exchange.Ledger[valuable.GetId()] = &ledgerEntry{
		ExchangeName: exchange.name,
		OpenShares:   amount,
		Holders:      make(map[string]float64),
	}
	return nil
}

// #########################
//			Portfolio
// #########################

func (exchange *Exchange) GetHoldingsCount(portfolio *portfolio.Portfolio, valuable valuable.Valuable) float64 {
	amount, exists := exchange.Ledger[valuable.GetId()].Holders[portfolio.UUID]
	if !exists {
		return 0
	}
	return amount
}

// #########################
// 			Trade
// #########################

func InitiateTrade(o *order.PurchaseOrder) {
	exchanger, exists := Exchanges[o.ExchangeID]
	if !exists {
		order.FailureOrder("exchange does not exist", o)
		return
	}
	exchanger.tradeChannel.Offer(o)
}

func (exchange *Exchange) StartExchange() {
	incomingTrades := exchange.tradeChannel.GetOutput()
	go func() {
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
	//get the stock if it exists
	value, exists := valuable.Valuables[o.ValuableID]
	if !exists {
		order.FailureOrder("asset is not recognized", o)
		return
	}
	value.GetLock().Acquire("trade")
	defer value.GetLock().Release()

	port, exists := portfolio.Portfolios[o.PortfolioID]
	if !exists {
		order.FailureOrder("portfolio does not exist, this is very bad", o)
		return
	}
	port.Lock.Acquire("trade")
	defer port.Lock.Release()

	info, ok := exchange.Ledger[value.GetId()]
	if !ok {
		order.FailureOrder("stock is not known to exchange", o)
		return
	}
	if o.Amount > 0 {
		//we have a buy
		// are there enough shares
		if info.OpenShares < o.Amount {
			order.FailureOrder("not enough shares", o)
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
		info.OpenShares -= o.Amount
		//add the holder amount
		_, exists := info.Holders[port.UUID]
		if !exists {
			info.Holders[port.UUID] = o.Amount
		} else {
			info.Holders[port.UUID] += o.Amount
		}
		// Update the portfolio with the new info
		port.TradeUpdate(value, info.Holders[port.UUID], costOfTrade)
		order.SuccessOrder(o)
	} else {
		// we have a sell
		//make sure they have that many shares
		currentOwns, exists := info.Holders[port.UUID]
		if !exists {
			order.FailureOrder("does not own any of stock", o)
			return
		}
		if currentOwns < o.Amount {
			order.FailureOrder("not enough shares", o)
			return
		}
		// make trade
		moneyReturned := o.Amount + value.GetValue()
		// add to open shares
		info.OpenShares += o.Amount
		// remove from ledger
		info.Holders[port.UUID] -= o.Amount
		if info.Holders[port.UUID] < TOLERANCE {
			delete(info.Holders, port.UUID)
			port.TradeUpdate(value, 0, moneyReturned)
		} else {
			port.TradeUpdate(value, info.Holders[port.UUID], moneyReturned)
		}
		order.SuccessOrder(o)
	}
	exchange.LedgerUpdateChannel.Offer(info)

}
