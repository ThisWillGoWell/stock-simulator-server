package messages

const TradeAction = "trade"
const ProspectTradeAction = "prospect"

type TradeMessage struct {
	StockId    string `json:"stock_id"`
	ExchangeID string `json:"-"`
	Amount     int64  `json:"amount"`
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
