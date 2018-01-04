package wallstreet

import "stock-server/utils"

type Portfolio struct {
	exchange *Exchange
	Wallet       float64
	NetWorth	 float64
	tradeChannel chan *PurchaseOrder

	//keeps track of how much $$$ they own, used for some slight optomization on calc networth
	personalLedger map[*Stock]float64

	updateChannel utils.ChannelDuplicator
	stockUpdates *utils.ChannelDuplicator
}

// this will set the networth value whenever the stock changes
func (port *Portfolio) stockUpdateListener(){
	stockUpdates := port.stockUpdates.GetOutput()
	// this will get called on each update
	// we are going to need to add all the stocks..
	// there is a race condition here, can you find it?.
	for stockUpdate := range stockUpdates{
		stock := stockUpdate.(*Stock)
		port.updateNetWorth(stock)
		port.updateChannel.Offer(port)
	}

}

func  (port *Portfolio) updateNetWorth(stock *Stock) {
	value := port.exchange.GetHoldingsCount(port, stock) * stock.CurrentPrice
	port.NetWorth = port.NetWorth - port.personalLedger[stock] + value
	port.personalLedger[stock] = value
}


func (port *Portfolio) tradeStock(ticker string, amount float64) {
	order := port.exchange.initiateTrade(ticker, port, amount)
	response := <- order.responseChannel
	if response.success{
		port.updateNetWorth()
	} else{
		//something went wrong
	}
}