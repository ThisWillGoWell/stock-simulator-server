package messages

const TradeAction = "trade"

type TradeMessage struct {
	StockTicker string  `json:"stock_ticker"`
	ExchangeID  string  `json:"exchange_id"`
	Amount      float64 `json:"amount"`
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

func BuildPurchaseResponse(message *TradeMessage, response interface{}) *BaseMessage {
	return &BaseMessage{
		Action:AlertAction,
		Msg: &AlertMessage{
			Alert: &TradeResponse{
				Trade:    message,
				Response: response,
			},
			Type: "trade_response",
		},
	}

}
