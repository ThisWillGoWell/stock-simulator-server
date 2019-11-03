package notification

import "fmt"

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

func DoneTradeNotification(portfilioUuid, stockUuid string, amount int64) error {
	_, err := NewNotification(portfilioUuid, TradeNotificationType, &TradeNotification{
		Amount:    amount,
		StockUuid: stockUuid,
		Success:   true,
	})
	return err
}

func SendMoneyTradeNotification(portfolioUuid, receiverUuid string, amount int64) error {
	_, err1 := NewNotification(portfolioUuid, SendMoneyNotificationType, MoneyTransferNotification{
		Sender:   portfolioUuid,
		Receiver: receiverUuid,
		Amount:   amount,
	})
	_, err2 := NewNotification(receiverUuid, RecieveNotificationType, MoneyTransferNotification{
		Sender:   portfolioUuid,
		Receiver: receiverUuid,
		Amount:   amount,
	})
	if err1 != nil || err2 != nil {
		return fmt.Errorf("send money transfer err1=[%v] err2=[%v]", err1, err2)
	}
	return nil
}
