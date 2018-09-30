package account

import (
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/messages"
	"github.com/stock-simulator-server/src/notification"
)

var GlobalNewObjects = duplicator.MakeDuplicator("global-new-objects")
var GlobalDeletes = duplicator.MakeDuplicator("global-deletes")
var GlobalNotifications = duplicator.MakeDuplicator("global-notifications")
var GlobalUpdates = duplicator.MakeDuplicator("global-new-objects")
var Globals = duplicator.MakeDuplicator("global-broadcast")

func RunGlobalSender() {
	go func() {
		out := GlobalNewObjects.GetBufferedOutput(10)
		for ele := range out {
			Globals.Offer(messages.NewObjectMessage(ele.(change.Identifiable)))
		}
	}()
	go func() {
		out := GlobalDeletes.GetBufferedOutput(10)
		for ele := range out {
			Globals.Offer(messages.BuildDeleteMessage(ele.(change.Identifiable)))
		}
	}()
	go func() {
		out := GlobalNotifications.GetBufferedOutput(10)
		for ele := range out {
			Globals.Offer(messages.BuildNotificationMessage(ele.(*notification.Notification)))
		}
	}()
	go func() {
		out := GlobalUpdates.GetBufferedOutput(10)
		for ele := range out {
			Globals.Offer(messages.BuildUpdateMessage(ele.(change.Identifiable)))
		}
	}()
}

type sender struct {
	lock          *lock.Lock
	activeClients int
	NewObjects    *duplicator.ChannelDuplicator
	Updates       *duplicator.ChannelDuplicator
	Notifications *duplicator.ChannelDuplicator
	Deletes       *duplicator.ChannelDuplicator
	Output        *duplicator.ChannelDuplicator
	close         chan interface{}
}

func newSender(userUuid string) *sender {

	sender := &sender{
		lock:          lock.NewLock("client-user-Sender-" + userUuid),
		activeClients: 0,
		NewObjects:    duplicator.MakeDuplicator("objects-Sender-" + userUuid),
		Updates:       duplicator.MakeDuplicator("update-Sender-" + userUuid),
		Deletes:       duplicator.MakeDuplicator("delete-Sender-" + userUuid),
		Notifications: duplicator.MakeDuplicator("notification-Sender-" + userUuid),
		Output:        duplicator.MakeDuplicator("output-Sender-" + userUuid),
		close:         make(chan interface{}),
	}
	sender.Output.RegisterInput(Globals.GetBufferedOutput(10))
	sender.runSendDeletes()
	sender.runSendObjects()
	sender.runSendUpdates()
	sender.runSendNotifications()
	return sender
}

func (s *sender) GetOutput() chan interface{} {
	s.lock.Acquire("new output")
	defer s.lock.Release()
	s.activeClients += 1
	return s.Output.GetBufferedOutput(10)
}

func (s *sender) CloseOutput(ch chan interface{}) {
	s.lock.Acquire("close output")
	defer s.lock.Release()
	s.Output.UnregisterOutput(ch)
	s.activeClients -= 1
}

func (s *sender) stop() {
	for i := 0; i < 4; i++ {
		s.close <- nil
	}
	duplicator.UnlinkDouplicator(s.Output, Globals)

	s.Notifications.StopDuplicator()
	s.Updates.StopDuplicator()
	s.NewObjects.StopDuplicator()
	s.Deletes.StopDuplicator()
	s.Output.StopDuplicator()

	close(s.close)
}

func (s *sender) runSendObjects() {
	object := s.NewObjects.GetBufferedOutput(10)
	go func() {
		for {
			select {
			case newObject := <-object:
				s.Output.Offer(messages.NewObjectMessage(newObject.(change.Identifiable)))
			case <-s.close:
				break
			}
		}
	}()
}

func (s *sender) runSendUpdates() {
	object := s.Updates.GetBufferedOutput(10)
	go func() {
		for {
			select {
			case newObject := <-object:
				s.Output.Offer(messages.BuildUpdateMessage(newObject.(change.Identifiable)))
			case <-s.close:
				break
			}
		}
	}()
}

func (s *sender) runSendDeletes() {
	object := s.Updates.GetBufferedOutput(10)
	go func() {
		for {
			select {
			case newObject := <-object:
				s.Output.Offer(messages.BuildUpdateMessage(newObject.(change.Identifiable)))
			case <-s.close:
				break
			}
		}
	}()
}

func (s *sender) runSendNotifications() {
	notifications := s.Notifications.GetBufferedOutput(10)
	go func() {
		for {
			select {
			case newNotifications := <-notifications:
				s.Output.Offer(messages.BuildNotificationMessage(newNotifications.(*notification.Notification)))
			case <-s.close:
				break
			}
		}
	}()
}
