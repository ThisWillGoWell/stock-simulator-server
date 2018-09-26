package transfer

import (
	"github.com/stock-simulator-server/src/order"
	"github.com/stock-simulator-server/src/portfolio"
)

func ExecuteTransfer(o *order.TransferOrder) {
	portfolio.PortfoliosLock.Acquire("moneyTransfer")
	defer portfolio.PortfoliosLock.Release()
	port, exists := portfolio.Portfolios[o.PortfolioID]
	if !exists {
		order.FailureOrder("giver portfolio not known", 0)
	}
	receiver, exists := portfolio.Portfolios[o.ReceiverID]
	if !exists {
		order.FailureOrder("portfolio does not exist, this is very bad", o)
		return
	}
	port.Lock.Acquire("transfer")
	defer port.Lock.Release()
	receiver.Lock.Acquire("transfer")
	defer receiver.Lock.Release()
	if o.Amount <= 0 {
		order.FailureOrder("invalid amount", o)
		return
	}
	if o.Amount > port.Wallet {
		order.FailureOrder("not enough money", o)
		return
	}
	receiver.Wallet += o.Amount
	port.Wallet -= o.Amount
	order.SuccessOrder(o)
	go port.Update()
	go receiver.Update()
}
