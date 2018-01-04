package wallstreet

import "stock-server/utils"

type Portfolio struct {
	stockManager StockManager
	worth float64

	updateChannel utils.ChannelDuplicator
	stockUpdates utils.ChannelDuplicator

}

func (port *Portfolio) listen(){
	go func() {
//		var stockUpdateChannel chan interface{}
//		stockUpdateChannel = port.stockUpdates.GetOutput()
//		for{
//			updatedObj:= <- stockUpdateChannel
//			updateStock := updatedObj.(Stock)
//			port.worth = updateStock.CurrentPrice
//			port.updateChannel.Offer(port)
//		}

		for stock := range port.stockUpdates.GetOutput(){
			port.worth = stock.(Stock).CurrentPrice
			port.updateChannel.Offer(port)
		}
	}()
}

func (port *Portfolio) buyStock(ticker string, amount float64) {
	stock := port.stockManager.getStock(ticker)

	port.worth  = stock.CurrentPrice
	port.stockUpdates.RegisterInput(stock.UpdateChannel.GetOutput())
	port.updateChannel.Offer(port)
}