package notification

import (
	"encoding/json"
	"time"

	"github.com/stock-simulator-server/src/change"

	"github.com/stock-simulator-server/src/sender"

	"github.com/stock-simulator-server/src/wires"

	"github.com/stock-simulator-server/src/lock"

	"github.com/stock-simulator-server/src/utils"

	"github.com/pkg/errors"
)

var notificationLock = lock.NewLock("notifications")
var notifications = make(map[string]*Notification)
var notificationsPortfolioUuid = make(map[string]map[string]*Notification)

const IdentifiableType = "notification"

type Notification struct {
	Uuid          string      `json:"uuid"`
	PortfolioUuid string      `json:"portfolio_uuid"`
	Timestamp     time.Time   `json:"time"`
	Type          string      `json:"type"`
	Notification  interface{} `json:"notification"`
	Seen          bool        `json:"seen"`
}

func NewNotification(portfolioUuid, t string, notification interface{}) *Notification {
	uuid := utils.SerialUuid()
	return MakeNotification(uuid, portfolioUuid, t, time.Now(), false, notification)
}

func MakeNotification(uuid, portfolioUuid, t string, timestamp time.Time, seen bool, notification interface{}) *Notification {
	notificationLock.Acquire("get-all-notifications")
	defer notificationLock.Release()
	note := &Notification{
		Uuid:          uuid,
		PortfolioUuid: portfolioUuid,
		Type:          t,
		Notification:  notification,
		Timestamp:     timestamp,
		Seen:          seen,
	}
	notifications[uuid] = note
	if _, ok := notificationsPortfolioUuid[portfolioUuid]; !ok {
		notificationsPortfolioUuid[portfolioUuid] = make(map[string]*Notification)
	}
	notificationsPortfolioUuid[portfolioUuid][uuid] = note
	utils.RegisterUuid(uuid, note)
	wires.NotificationNewObject.Offer(note)
	sender.SendNewObject(portfolioUuid, note)
	return note
}

func AcknowledgeNotification(uuid, portfolioUuid string) error {
	notification, ok := notifications[uuid]
	if !ok {
		return errors.New("notification uuid does not exist")
	}
	if notification.PortfolioUuid != portfolioUuid {
		return errors.New("user does not own notification, what are you doing?")
	}
	notification.Seen = true
	sender.SendChangeUpdate(notification.PortfolioUuid, &change.ChangeNotify{
		Type:    notification.GetType(),
		Id:      notification.GetId(),
		Object:  notification,
		Changes: []*change.ChangeField{{Field: "seen", Value: true}},
	})
	return nil
}

type MailNotification struct {
	From  string `json:"from"`
	Text  string `json:"text"`
	Money int64  `json:"money"`
}

func NewMailNotifcation(uuid, from string, text string, money int64) *Notification {
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

func GetAllNotifications(portfolioUuid string) []*Notification {
	notificationLock.Acquire("get-all-notifications")
	defer notificationLock.Release()
	notifications := make([]*Notification, 0)
	for _, notification := range notificationsPortfolioUuid[portfolioUuid] {
		notifications = append(notifications, notification)
	}
	return notifications
}

func JsonToNotification(jsonString, notificationType string) interface{} {
	var i interface{}
	switch notificationType {
	case NewItemNotificationType:
		i = ItemNotification{}
	case UsedItemNotificationType:
		i = ItemNotification{}
	case TradeNotificationType:
		i = TradeNotification{}
	case SendMoneyNotificationType:
		i = MoneyTransferNotification{}
	case RecieveNotificationType:
		i = MoneyTransferNotification{}
	case NewEffectNotificationType:
		i = EffectNotification{}
	case EndEffectNotificationType:
		i = EffectNotification{}
	}

	json.Unmarshal([]byte(jsonString), &i)
	return &i
}

func DeleteNotification(uuid, portfolioUuid string) error {
	notificationLock.Acquire("delete note")
	defer notificationLock.Release()
	if _, exists := notificationsPortfolioUuid[portfolioUuid]; !exists {
		return errors.New("user does not have any notification")
	}
	note, exists := notificationsPortfolioUuid[portfolioUuid][uuid]
	if !exists {
		return errors.New("notification does not exist")
	}
	delete(notifications, uuid)
	delete(notificationsPortfolioUuid[note.PortfolioUuid], uuid)
	if len(notificationsPortfolioUuid[note.PortfolioUuid]) == 0 {
		delete(notificationsPortfolioUuid, note.PortfolioUuid)
	}

	utils.RemoveUuid(uuid)
	sender.SendDeleteObject(portfolioUuid, note)
	wires.NotificationsDelete.Offer(note)
	return nil
}

func (note *Notification) GetId() string {
	return note.Uuid
}

func (*Notification) GetType() string {
	return IdentifiableType
}

//**
// todo
//
//func StartCleanNotifications() {
//	go runCleanNotifications()
//}
//func runCleanNotifications() {
//	for {
//		userListLock.Acquire("clean notifications")
//		for _, user := range UserList {
//			user.Lock.Acquire("clean notifications")
//			newStartIndex := 0
//			for _, notification := range user.Notifications {
//				if time.Since(notification.Timestamp) < notificationTimeLimit {
//					break
//				} else {
//					newStartIndex += 1
//				}
//			}
//			user.Notifications = user.Notifications[newStartIndex:]
//			userListLock.Release()
//		}
//		userListLock.Release()
//		<-time.After(time.Hour)
//	}
//}
