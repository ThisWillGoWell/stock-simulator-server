package notification

const TradeNotificationType = "trade"
const SendMoneyNotificationType = "send_money"
const ReceiveNotificationType = "receive_money"

type TradeNotification struct {
	Success   bool   `json:"success"`
	Amount    int64  `json:"amount"`
	StockUuid string `json:"stock"`
	Err       string `json:"error,omitempty"`
}

type MoneyTransferNotification struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int64  `json:"amount"`
}

func DoneTradeNotification(portfilioUuid, stockUuid string, amount int64) *Notification {
	return NewNotification(portfilioUuid, TradeNotificationType, &TradeNotification{
		Amount:    amount,
		StockUuid: stockUuid,
		Success:   true,
	})
}

func SendMoneyTradeNotification(portfolioUuid, receiverUuid string, amount int64) (*Notification, *Notification) {
	n1 := NewNotification(portfolioUuid, SendMoneyNotificationType, MoneyTransferNotification{
		Sender:   portfolioUuid,
		Receiver: receiverUuid,
		Amount:   amount,
	})
	n2 := NewNotification(receiverUuid, ReceiveNotificationType, MoneyTransferNotification{
		Sender:   portfolioUuid,
		Receiver: receiverUuid,
		Amount:   amount,
	})

	return n1, n2
}
