package order

import (
	"github.com/stock-simulator-server/src/trade"
	"github.com/stock-simulator-server/src/transfer"
)

type Order interface {
	rc() chan *Response
}

var orderChannel = make(chan Order, 30)



func Run() {
	go func() {
		for o := range orderChannel {
			switch o.(type) {
			case *PurchaseOrder:
				trade.ExecuteTrade(o.(*PurchaseOrder))
			case *TransferOrder:
				transfer.ExecuteTransfer(o.(*TransferOrder))
			}
		}
	}()
}


type TransferOrder struct{
	ReceiverID string `json:"receiver"`
	PortfolioID string 	`json:"giver"`
	Amount int64 `json:"amount"`
	ResponseChannel chan *Response
}
func (to *TransferOrder) rc()chan *Response{return to.ResponseChannel}


/**
Purchase order represents a single order
*/
type PurchaseOrder struct {
	ValuableID      string         `json:"stock_id"`
	PortfolioID     string         `json:"portfolio"`
	ExchangeID      string         `json:"exchange"`
	Amount          int64          `json:"amount"`
	ResponseChannel chan *Response `json:"-"`
}
func (po *PurchaseOrder) rc() chan *Response {return po.ResponseChannel}


type Response struct {
	Order   Order 			`json:"order"`
	Success bool           `json:"success"`
	Err     string         `json:"err"`
}

// note this does not validate if the stock exists or not, that's done in the trade() function
func MakePurchaseOrder(valuableID, portfolioUUID string, amount int64) *PurchaseOrder {
	po := &PurchaseOrder{
		//ExchangeID:      exchangeID,
		ValuableID:      valuableID,
		PortfolioID:     portfolioUUID,
		Amount:          amount,
		ResponseChannel: make(chan *Response, 1),
	}
	orderChannel <- po
	return po
}

func MakeTransferOrder(giver, reciever string, amount int64) *TransferOrder{
	to := &TransferOrder{
		PortfolioID: giver,
		ReceiverID: reciever,
		Amount: amount,
	}
	orderChannel <- to
	return to
}

func (po *PurchaseOrder) Execute() {

}

func SuccessOrder(o Order) {
	o.rc() <- &Response{
		Order:   o,
		Success: true,
	}
}
func FailureOrder(msg string, o Order) {

	o.rc() <- &Response{
		Order:   o,
		Success: false,
		Err:     msg,
	}
}
