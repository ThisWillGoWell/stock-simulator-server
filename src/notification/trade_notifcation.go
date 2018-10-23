package notification

const TradeNotificationType = "trade"
const SendMoneyNotificationType = "send_money"
const RecieveNotificationType = "receive_money"

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

func DoneTradeNotification(userUuid, stockUuid string, amount int64) {
	NewNotification(userUuid, TradeNotificationType, &TradeNotification{
		Amount:    amount,
		StockUuid: stockUuid,
		Success:   true,
	})
}

func SendMoneyTradeNotification(senderUuid, receiverUuid string, amount int64) {
	NewNotification(senderUuid, SendMoneyNotificationType, MoneyTransferNotification{
		Sender:   senderUuid,
		Receiver: receiverUuid,
		Amount:   amount,
	})
	NewNotification(receiverUuid, RecieveNotificationType, MoneyTransferNotification{
		Sender:   senderUuid,
		Receiver: receiverUuid,
		Amount:   amount,
	})
}
