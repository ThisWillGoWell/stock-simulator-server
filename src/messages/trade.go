package messages

const TradeAction = "trade"

type TradeMessage struct {
	StockId    string  `json:"stock_id"`
	ExchangeID string  `json:"exchange_id"`
	Amount     int64 `json:"amount"`
}

func (*TradeMessage) message() { return }

type TradeResponse struct {
	Trade    *TradeMessage `json:"trade"`
	Response interface{}   `json:"response"`
}

func (*TradeResponse) message() { return }

func (baseMessage *BaseMessage) IsTrade() bool {
	return baseMessage.Action == "trade"
}

func BuildPurchaseResponse(response interface{}) *BaseMessage {
	return &BaseMessage{
		Action: AlertAction,
		Msg: &AlertMessage{
			Alert: response,
			Type:  "trade_response",
		},
	}

}
