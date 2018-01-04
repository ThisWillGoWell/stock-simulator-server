package wallstreet

import (
	"time"
	"plugin"
	"fmt"
	"stock-server/utils"
)

type Exchange struct {
	stockManager *StockManager
	ledger map[*Stock]*ledgerEntry
	tradeChannel *utils.ChannelDuplicator
	updateIntervalLength time.Duration

}

type ledgerEntry struct {
	stock      *Stock
	holders    map[*Portfolio]float64
	openShares float64
}


func BuildExchange(interval time.Duration) *Exchange{
	return &Exchange{
		stockManager:         buildStockManager(),
		ledger:               make(map[*Stock]*ledgerEntry),
		tradeChannel:         utils.MakeDuplicator(),
		updateIntervalLength: interval,
	}
}
// #########################
//			Stocks
// #########################

// calculate the percent of a stock to run on a given delta
func  (exchange *Exchange) calculateRunPercent(targetInterval time.Duration) float64 {
	return targetInterval.Seconds() / exchange.updateIntervalLength.Seconds()
}

func (exchange *Exchange)AddStock(tickerId string, name string, startPrice, totalShares float64, duration time.Duration){
	stock := exchange.stockManager.addStock(tickerId, name, startPrice, exchange.calculateRunPercent(duration))
	if stock == nil {
		// something went wrong
		return
	}
	exchange.ledger[stock] = buildLedgerEntry(stock, totalShares)
}

func buildLedgerEntry(stock *Stock, totalShares float64) *ledgerEntry{
	return &ledgerEntry{
		stock: stock,
		openShares: totalShares,
		holders: make( map[*Portfolio]float64),
	}
}
// #########################
//			Portfolio
// #########################

func (exchange *Exchange)AddPortfolio(uuid string, ){

}

func (exchange *Exchange)GetHoldingsCount(portfolio *Portfolio, stock *Stock) float64{
	return exchange.ledger[stock].holders[portfolio]
}


// #########################
// 			Trade
// #########################

type PurchaseOrder struct {
	stock           *Stock
	portfolio       *Portfolio
	amount          float64
	responseChannel chan *PurchasedResponse
}


type PurchasedResponse struct {
	success bool
	err error
}

type ProcessError struct {
	order *PurchaseOrder
	message string
}
// note this does not validate if the stock exists or not, thats done in the trade() funciton
func (exchange *Exchange )buildPurchaseOrder(stock *Stock, portfolio *Portfolio, amount float64) *PurchaseOrder {
	return &PurchaseOrder{
		stock:           stock,
		portfolio:       portfolio,
		amount:          amount,
		responseChannel: make(chan *PurchasedResponse),
	}
}

func (exchange *Exchange) initiateTrade(stockTicker string, portfolio *Portfolio, amount float64) *PurchaseOrder{
	stock := exchange.stockManager.getStock(stockTicker)
	order := exchange.buildPurchaseOrder(stock, portfolio, amount)
	exchange.tradeChannel.Offer(order)
	return order
}

func (e *ProcessError) Error() string {
	return e.message
}

func (exchange *Exchange)runExchange(){
	deltaTicker := time.NewTicker(exchange.updateIntervalLength)
	incomingTrades := exchange.tradeChannel.GetOutput()
	for {
		select{
		case purchaseOrder := <-incomingTrades:
			exchange.trade(purchaseOrder.(*PurchaseOrder))
		case <-deltaTicker.C:
			exchange.stockManager.changeStock()

		}
	}
}

func failureTrade(msg string, order *PurchaseOrder){
	order.responseChannel <- &PurchasedResponse{
		success: false,
		err: &ProcessError{
			message: msg,
			order: order,
		},
	}
}

func successPurchase(order *PurchaseOrder){
	order.responseChannel <- &PurchasedResponse{
		success:true,
		err: nil,
	}
}

// validate and make a trade
// Purchase Order contains a reference to the order
// use failureTrade() and successPurchase to send response down channel
func (exchange *Exchange) trade(order *PurchaseOrder){
	stockInfo, ok := exchange.ledger[order.stock]
	if !ok {
		failureTrade("stock is not known", order)
	}
	if order.amount > 0{
		//we have a buy
		// are there enough shares
		if stockInfo.openShares >= order.amount{
			failureTrade("not enough shares", order)
			return
		}
		// does the user have enough money
		costOfTrade := order.amount * stockInfo.stock.CurrentPrice
		if costOfTrade > order.portfolio.Wallet {
			failureTrade("not enough $$", order)
			return
		}


		// any other checks?

		// make the trade
		// subtract from open shares
		stockInfo.openShares -= order.amount
		//add the holder amount
		stockInfo.holders[order.portfolio] += order.amount
		// subtract from the wallet
		order.portfolio.Wallet -= costOfTrade
		successPurchase(order)
	} else {
		// we have a sell
		//make sure they have that many shares
		currentOwns := stockInfo.holders[order.portfolio]
		if currentOwns < order.amount{
			failureTrade("not enough shares", order)
		}
		// make trade
		// add the cash into the wallet
		order.portfolio.Wallet += order.amount + stockInfo.stock.CurrentPrice
		// add to open shares
		exchange.ledger[order.stock].openShares += order.amount
		// remove from ledger
		exchange.ledger[order.stock].holders[order.portfolio] -= order.amount
		successPurchase(order)
	}
}
