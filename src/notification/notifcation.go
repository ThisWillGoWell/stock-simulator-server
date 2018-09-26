package notification

import (
	"time"
)

type Notification struct{
	Timestamp time.Time `json:"time"`
	Type string `json:"type"`
	Notification interface{} `json:"notification"`
}


type TradeNotification struct{
	Success bool `json:"success"`
	Amount int64 `json:"amount"`
	StockUuid string `json:"stock"`
	Err		string `json:"error,omitempty"`
}

func NewTradeNotifcation(success bool, amount int64, stockUuid string, err error) *Notification{
	return &Notification{
		Timestamp: time.Now(),
		Type:      "trade",
		Notification: &TradeNotification{
			Success:   success,
			Amount:    amount,
			StockUuid: stockUuid,
			Err:       err.Error(),
		},
	}
}

type MailNotification struct {
	From string `json:"from"`
	Text string `json:"text"`
	Money int64 `json:"money"`
}

func NewMailNotifcation(from string, text string, money int64) *Notification{
	return &Notification{
		Timestamp: time.Now(),
		Type:      "mail",
		Notification: &MailNotification{
			From:  from,
			Text:  text,
			Money: money,
		},
	}
}