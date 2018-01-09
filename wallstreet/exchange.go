package wallstreet

import (
	"stock-server/utils"
	"time"
)

var exchanges map[string]*Exchange



const TOLERANCE = 0.000001

type Exchange struct {
	stockManager *StockManager
	ledger map[*Stock]*ledgerEntry
	tradeChannel *utils.ChannelDuplicator
	portfolios map[*Portfolio]bool
	updateIntervalLength time.Duration
	lock *utils.Lock
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
		lock: utils.NewLock(),
	}
}
// #########################
//			Stocks
// #########################

// calculate the percent of a stock to run on a given delta
func  (exchagne *Exchange) calculateRunPercent(targetInterval time.Duration) float64 {
	return targetInterval.Seconds() / exchagne.updateIntervalLength.Seconds()
}

func (exchagne *Exchange)AddStock(tickerId string, name string, startPrice, totalShares float64, duration time.Duration){
	exchagne.lock.Acquire()
	defer exchagne.lock.Release()

	stock := exchagne.stockManager.addStock(tickerId, name, startPrice, exchagne.calculateRunPercent(duration))
	if stock == nil {
		// something went wrong
		return
	}
	exchagne.ledger[stock] = buildLedgerEntry(stock, totalShares)
}

func buildLedgerEntry(stock *Stock, totalShares float64) *ledgerEntry{
	return &ledgerEntry{
		stock: stock,
		openShares: totalShares,
		holders: make( map[*Portfolio]float64),
	}
}

func (exchagne *Exchange)GetStockUpdateChanel() chan interface{}{
	return exchagne.stockManager.StockUpdateChannel.GetOutput()
}
// #########################
//			Portfolio
// #########################

func (exchagne *Exchange)GetHoldingsCount(portfolio *Portfolio, stock *Stock) float64{
	amount, exists := exchagne.ledger[stock].holders[portfolio]
	if !exists{
		return 0
	}
	return amount
}


// #########################
// 			Trade
// #########################

type PurchaseOrder struct {
	stock           *Stock
	seller          *Portfolio
	buyer			*Portfolio
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
func (exchagne *Exchange )buildPurchaseOrder(stock *Stock, portfolio *Portfolio, amount float64) *PurchaseOrder {
	return &PurchaseOrder{
		stock:           stock,
		seller:          portfolio,
		amount:          amount,
		responseChannel: make(chan *PurchasedResponse),
	}
}

func (exchagne *Exchange) initiateTrade(stockTicker string, portfolio *Portfolio, amount float64) *PurchaseOrder{
	stock := exchagne.stockManager.getStock(stockTicker)
	order := exchagne.buildPurchaseOrder(stock, portfolio, amount)
	exchagne.tradeChannel.Offer(order)
	return order
}

func (e *ProcessError) Error() string {
	return e.message
}


func (exchagne *Exchange)RunExchange(){
	deltaTicker := time.NewTicker(exchagne.updateIntervalLength)
	incomingTrades := exchagne.tradeChannel.GetOutput()
	for {
		select{
		case purchaseOrder := <-incomingTrades:
			exchagne.trade(purchaseOrder.(*PurchaseOrder))
		case <-deltaTicker.C:
			exchagne.stockManager.changeStock()
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

func successTrade(order *PurchaseOrder){
	order.responseChannel <- &PurchasedResponse{
		success:true,
		err: nil,
	}
}

// validate and make a trade
// Purchase Order contains a reference to the order
// use failureTrade() and successTrade to send response down channel
// Don't need to a lock around this since the portfolio holds it for that trade
func (exchagne *Exchange) trade(order *PurchaseOrder) {
	order.stock.lock.Acquire()
	defer order.stock.lock.Release()
	exchagne.tradeWithFloor(order)
}



func (exchagne *Exchange) tradeWithFloor(order *PurchaseOrder){
	stockInfo, ok := exchagne.ledger[order.stock]
	if !ok {
		failureTrade("stock is not known", order)
		return
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
		if costOfTrade > order.seller.Wallet {
			failureTrade("not enough $$", order)
			return
		}
		// any other checks?

		// make the trade
		// subtract from open shares
		stockInfo.openShares -= order.amount
		//add the holder amount
		_, exists := stockInfo.holders[order.seller]
		if ! exists {
			stockInfo.holders[order.seller] = order.amount
		}else{
			stockInfo.holders[order.seller] += order.amount
		}
		// subtract from the wallet
		order.seller.Wallet -= -1 * costOfTrade
		successTrade(order)
	} else {
		// we have a sell
		//make sure they have that many shares
		currentOwns, exists := stockInfo.holders[order.seller]
		if ! exists{
			failureTrade("does not own any of stock", order)
		}
		if currentOwns < order.amount{
			failureTrade("not enough shares", order)
		}
		// make trade
		// add the cash into the wallet
		order.seller.Wallet += order.amount + stockInfo.stock.CurrentPrice
		// add to open shares
		stockInfo.openShares += order.amount
		// remove from ledger
		stockInfo.holders[order.seller] -= order.amount
		if stockInfo.holders[order.seller] < TOLERANCE {
			delete(stockInfo.holders, order.seller)
		}
		successTrade(order)
	}
}

func (exchagne *Exchange) tradeWithTwo(order *PurchaseOrder) {
	order.buyer.lock.Acquire()
	defer order.buyer.lock.Release()

	stockInfo, ok := exchagne.ledger[order.stock]
	if !ok {
		failureTrade("stock is not known", order)
		return
	}
	//make sure we are only doing a positive trade
	if order.amount < 0 {
		failureTrade("order amount must be postitve", order)
		return
	}
	//make sure buyer has enough $$$
	costOfTrade := order.stock.CurrentPrice * order.amount
	if order.buyer.Wallet < costOfTrade{
		failureTrade("buyer does not have enough $$$", order)
		return
	}

	// make sure the seller has enough stocks
	sellerAmount, ok := stockInfo.holders[order.seller]
	if !ok {
		failureTrade("seller does not have that stock", order)
		return
	}

	if sellerAmount < order.amount {
		failureTrade("seller does not have enough stocks", order)
		return
	}
	//make trade
	stockInfo.holders[order.seller] -= order.amount
	stockInfo.holders[order.buyer] += order.amount
	order.buyer.Wallet -= costOfTrade
	order.seller.Wallet += costOfTrade

	successTrade(order)
	return
}
