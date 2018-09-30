package messages

import (
	"time"

	"github.com/stock-simulator-server/src/notification"
)

const NotificationAction = "notification"
const NotificationAck = "ack"

func (baseMessage *BaseMessage) IsNotification() bool {
	return baseMessage.Action == NotificationAction
}

type NotificationMessage struct {
	Notification interface{} `json:"notification"`
	Type         string      `json:"type"`
	Timestamp    time.Time   `json:"timestamp"`
	Seen         bool        `json:"seen"`
	Uuid         string      `json:"uuid"`
}

func (*NotificationMessage) message() { return }

func (baseMessage *BaseMessage) IsNotificationAck() bool {
	return baseMessage.Action == NotificationAck
}

type NotificationAckMessage struct {
	Uuid string `json:"uuid"`
}

func (*NotificationAckMessage) message() { return }

func BuildNotificationMessage(n *notification.Notification) *BaseMessage {
	return &BaseMessage{
		Action: NotificationAction,
		Msg: &NotificationMessage{
			Notification: n.Notification,
			Seen:         n.Seen,
			Uuid:         n.Uuid,
			Timestamp:    n.Timestamp,
			Type:         n.Type,
		},
	}
}

func NewErrorMessage(err string) *BaseMessage {
	return &BaseMessage{
		Action: NotificationAction,
		Msg: &NotificationMessage{
			Type:         "error",
			Notification: err,
			Timestamp:    time.Now(),
			Seen:         false,
		},
	}
}
