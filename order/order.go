package order

import (
	"errors"
)

type PurchaseOrder struct {
	ValuableID       string
	PortfolioID      string
	ExchangeID		string
	Amount          float64
	ResponseChannel chan *PurchasedResponse
}

type PurchasedResponse struct {
	Success bool
	Err error
}


// note this does not validate if the stock exists or not, thats done in the trade() funciton
func BuildPurchaseOrder(valuableID, exchangeID,portfolioUUID string, amount float64) *PurchaseOrder {
	return &PurchaseOrder{
		ExchangeID: 		exchangeID,
		ValuableID:        valuableID,
		PortfolioID:       portfolioUUID,
		Amount:          amount,
		ResponseChannel: make(chan *PurchasedResponse, 1),
	}
}

func (o *PurchaseOrder)Execute(){

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
