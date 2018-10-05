package notification

import (
	"encoding/json"
	"time"

	"github.com/stock-simulator-server/src/wires"

	"github.com/stock-simulator-server/src/account"

	"github.com/stock-simulator-server/src/lock"

	"github.com/stock-simulator-server/src/utils"

	"github.com/pkg/errors"
)

var notifcationLock = lock.NewLock("notifications")
var notifications = make(map[string]*Notification)
var notificationsUserUuid = make(map[string]map[string]*Notification)

type Notification struct {
	Uuid         string      `json:"uuid"`
	UserUuid     string      `json:"user_uuid"`
	Timestamp    time.Time   `json:"time"`
	Type         string      `json:"type"`
	Notification interface{} `json:"notification"`
	Seen         bool        `json:"seen"`
}

func NewNotification(userUuid, t string, notification interface{}) *Notification {
	uuid := utils.SerialUuid()
	return MakeNotification(uuid, userUuid, t, time.Now(), false, notification)
}

func MakeNotification(uuid, userUuid, t string, timestamp time.Time, seen bool, notification interface{}) *Notification {
	notifcationLock.Acquire("get-all-notifications")
	defer notifcationLock.Release()
	note := &Notification{
		Uuid:         uuid,
		UserUuid:     userUuid,
		Type:         t,
		Notification: notification,
		Timestamp:    timestamp,
		Seen:         seen,
	}
	notifications[uuid] = note
	if _, ok := notificationsUserUuid[userUuid]; !ok {
		notificationsUserUuid[userUuid] = make(map[string]*Notification)
	}
	notificationsUserUuid[userUuid][uuid] = note
	utils.RegisterUuid(uuid, note)
	wires.NotificationNewObject.Offer(note)
	account.SendNotifcation(userUuid, note)
	return note
}

func AcknowledgeNotification(uuid, userUuid string) error {
	notification, ok := notifications[uuid]
	if !ok {
		return errors.New("notification uuid does not exist")
	}
	if notification.Uuid != userUuid {
		return errors.New("user does not own notification, what are you doing?")
	}
	notification.Seen = true
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

func GetAllNotifications(userUuid string) []*Notification {
	notifcationLock.Acquire("get-all-notifications")
	defer notifcationLock.Release()
	notifications := make([]*Notification, 0)
	for _, notification := range notificationsUserUuid[userUuid] {
		notifications = append(notifications, notification)
	}
	return notifications
}

func JsonToNotifcation(jsonString, notifactionType string) interface{} {
	var i interface{}
	switch notifactionType {
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
	}
	json.Unmarshal([]byte(jsonString), &i)
	return i
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
