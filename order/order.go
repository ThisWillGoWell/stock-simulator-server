package order

import (
	"github.com/stock-simulator-server/valuable"
	"github.com/stock-simulator-server/portfolio"
	"errors"
)

type PurchaseOrder struct {
	Valuable        valuable.Valuable
	Portfolio       *portfolio.Portfolio
	Amount          float64
	ResponseChannel chan *PurchasedResponse
}

type PurchasedResponse struct {
	Success bool
	Err error
}


// note this does not validate if the stock exists or not, thats done in the trade() funciton
func BuildPurchaseOrder(valuable valuable.Valuable, portfolio *portfolio.Portfolio, amount float64) *PurchaseOrder {
	return &PurchaseOrder{
		Valuable:         valuable,
		Portfolio:          portfolio,
		Amount:          amount,
		ResponseChannel: make(chan *PurchasedResponse, 1),
	}
}

func SuccessOrder(o *PurchaseOrder){
	o.ResponseChannel <- &PurchasedResponse{
		Success:true,
		Err: nil,
	}
}
func FailureOrder(msg string, o *PurchaseOrder){
	o.ResponseChannel <- &PurchasedResponse{
		Success: false,
		Err:errors.New(msg),
	}
}
