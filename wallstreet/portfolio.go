package wallstreet

import "stock-server/utils"

var PortfolioIds map[string]*Portfolio


type Portfolio struct {
	Name     string
	UUID     string
	Exchange map[string]Exchange
	Wallet   float64
	NetWorth float64
	//keeps track of how much $$$ they own, used for some slight optomization on calc networth
	personalLedger map[*Stock]float64


	UpdateChannel *utils.ChannelDuplicator
	stockUpdates  *utils.ChannelDuplicator

	lock *utils.Lock
}

func NewPortfolio( userUUID,  name string )(*Portfolio){
	return &Portfolio{
		Name:           name,
		UUID:           userUUID,
		Wallet:         1000,
		NetWorth:       1000,
		personalLedger: make(map[*Stock]float64),
		UpdateChannel:  utils.MakeDuplicator(),
		stockUpdates:   utils.MakeDuplicator(),
		lock :          utils.NewLock(),
	}
}

func (port *Portfolio) JoinExchange(name string, exchange *Exchange){
	port.Exchange = append(port.Exchange, exchange)
}
// this will set the networth value whenever the stock changes
func (port *Portfolio) stockUpdateListener(){
	stockUpdates := port.stockUpdates.GetOutput()
	// this will get called on each update
	// we are going to need to add all the stocks..
	// there is a race condition here, can you find it?.
	for stockUpdate := range stockUpdates{
		port.lock.Acquire()
		stock := stockUpdate.(*Stock)
		port.updateNetWorth(stock)
		port.lock.Release()
	}
}

func  (port *Portfolio) updateNetWorth(stock *Stock) {
	value := port.Exchange.GetHoldingsCount(port, stock) * stock.CurrentPrice
	port.NetWorth = port.NetWorth - port.personalLedger[stock] + value
	port.personalLedger[stock] = value
	port.UpdateChannel.Offer(port)

}

func (port *Portfolio) updateWallet(walletChange float64){
	port.NetWorth += walletChange
}

func (port *Portfolio) tradeStock(ticker string, amount float64) {
	port.lock.Acquire()
	defer port.lock.Release()
	order := port.Exchange.initiateTrade(ticker, port, amount)
	response := <- order.responseChannel
	if response.success{
		port.updateNetWorth(order.stock)
	} else{
		//something went wrong
	}

}