package notification

import (
	"encoding/json"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id/change"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/sender"

	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"

	"fmt"
	"github.com/ThisWillGoWell/stock-simulator-server/src/database"
	"github.com/pkg/errors"
)

var NotificationLock = lock.NewLock("notifications")
var notifications = make(map[string]*Notification)
var notificationsPortfolioUuid = make(map[string]map[string]*Notification)

const IdentifiableType = "notification"

type Notification struct {
	objects.Notification
}

func NewNotification(portfolioUuid, t string, notification interface{}) *Notification {
	uuid := id.SerialUuid()

	noti :=  objects.Notification{
		Uuid:          uuid,
		PortfolioUuid: portfolioUuid,
		Type:          t,
		Notification:  notification,
		Timestamp:     time.Now(),
		Seen:          false,
	}

	n := MakeNotification(noti)
	return n
}

func DeleteNotification(uuid string, lockAcquired bool) error{
	if !lockAcquired {
		NotificationLock.Acquire("delete note")
		defer NotificationLock.Release()
	}

	n, ok := notifications[uuid]
	if !ok {
		log.Log.Errorf("got a delete for a uuid that does not exists")
		return nil
	}
	if dbErr := database.Db.Execute(nil, []interface{}{n}); dbErr != nil {
		log.Log.Errorf("failed to delete notification from database err=[%v]", dbErr)
		return fmt.Errorf("opps! something went wrong 0x231")
	}
	deleteNotification(n)
	sender.SendDeleteObject(n.Uuid, n)
	return nil
}

func deleteNotification(note *Notification){
	delete(notifications, note.Uuid)
	if _, exists := notificationsPortfolioUuid[note.PortfolioUuid]; !exists {
		log.Log.Errorf("user does not have any notifications ")
		return
	}
	note, exists := notificationsPortfolioUuid[note.PortfolioUuid][note.Uuid]
	if !exists {
		log.Log.Errorf("notification does not exist in users inventory")
		return
	}
	delete(notificationsPortfolioUuid[note.PortfolioUuid], note.Uuid)
	if len(notificationsPortfolioUuid[note.PortfolioUuid]) == 0 {
		delete(notificationsPortfolioUuid, note.PortfolioUuid)
	}
	id.RemoveUuid(note.Uuid)

}

func MakeNotification(notification objects.Notification) *Notification {
	if s, ok := notification.Notification.(string); ok {
		notification.Notification = JsonToNotification(s, notification.Type)
	}
	note := &Notification{
		notification,
	}
	notifications[notification.Uuid] = note
	if _, ok := notificationsPortfolioUuid[notification.PortfolioUuid]; !ok {
		notificationsPortfolioUuid[notification.PortfolioUuid] = make(map[string]*Notification)
	}
	notificationsPortfolioUuid[notification.PortfolioUuid][note.Uuid] = note
	id.RegisterUuid(note.Uuid, note)
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
		objects.Notification{
			Timestamp: time.Now(),
			Type:      "mail",
			Notification: &MailNotification{
				From:  from,
				Text:  text,
				Money: money,
			},
		}}
}

func GetAllNotifications(portfolioUuid string) []*Notification {
	NotificationLock.Acquire("get-all-notifications")
	defer NotificationLock.Release()
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
	case ReceiveNotificationType:
		i = MoneyTransferNotification{}
	case NewEffectNotificationType:
		i = EffectNotification{}
	case EndEffectNotificationType:
		i = EffectNotification{}
	}

	json.Unmarshal([]byte(jsonString), &i)
	return &i
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
