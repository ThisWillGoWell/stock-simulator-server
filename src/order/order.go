package order

import (
	"errors"
)

type PurchaseOrder struct {
	ValuableID      string `json:"asset"`
	PortfolioID     string `json:"portfolio"`
	ExchangeID      string `json:"exchange"`
	Amount          float64 `json:"amount"`
	ResponseChannel chan *PurchasedResponse `json:"-"`
}

type PurchasedResponse struct {
	Order *PurchaseOrder `json:"order"`
	Success bool `json:"success"`
	Err     error `json:"err"`
}

// note this does not validate if the stock exists or not, thats done in the trade() funciton
func BuildPurchaseOrder(valuableID, exchangeID, portfolioUUID string, amount float64) *PurchaseOrder {
	return &PurchaseOrder{
		ExchangeID:      exchangeID,
		ValuableID:      valuableID,
		PortfolioID:     portfolioUUID,
		Amount:          amount,
		ResponseChannel: make(chan *PurchasedResponse, 1),
	}
}

func (o *PurchaseOrder) Execute() {

}

func SuccessOrder(o *PurchaseOrder) {
	o.ResponseChannel <- &PurchasedResponse{
		Order: o,
		Success: true,
		Err:     nil,
	}
}
func FailureOrder(msg string, o *PurchaseOrder) {
	o.ResponseChannel <- &PurchasedResponse{
		Order: o,
		Success: false,
		Err:     errors.New(msg),
	}
}
