package sender

import (
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/messages"
	"github.com/stock-simulator-server/src/wires"
)

var Senders = make(map[string]*Sender)

var GlobalMessages = duplicator.MakeDuplicator("global-messages")

func RunGlobalSender() {
	go func() {
		out := wires.GlobalNewObjects.GetBufferedOutput(10000)
		for ele := range out {
			GlobalMessages.Offer(messages.NewObjectMessage(ele.(change.Identifiable)))
		}
	}()
	go func() {
		out := wires.GlobalDeletes.GetBufferedOutput(10000)
		for ele := range out {
			GlobalMessages.Offer(messages.BuildDeleteMessage(ele.(change.Identifiable)))
		}
	}()
	go func() {
		out := wires.GlobalNotifications.GetBufferedOutput(10000)
		for ele := range out {
			GlobalMessages.Offer(messages.BuildNotificationMessage(ele))
		}
	}()
	go func() {
		out := change.PublicSubscribeChange.GetBufferedOutput(100000)
		for ele := range out {
			GlobalMessages.Offer(messages.BuildUpdateMessage(ele.(change.Identifiable)))
		}
	}()
}

type Sender struct {
	lock          *lock.Lock
	activeClients int
	NewObjects    *duplicator.ChannelDuplicator
	Updates       *duplicator.ChannelDuplicator
	Notifications *duplicator.ChannelDuplicator
	Deletes       *duplicator.ChannelDuplicator
	Output        *duplicator.ChannelDuplicator
	close         chan interface{}
}

func NewSender(userUuid string) *Sender {
	sender := &Sender{
		lock:          lock.NewLock("client-user-Sender-" + userUuid),
		activeClients: 0,
		NewObjects:    duplicator.MakeDuplicator("objects-Sender-" + userUuid),
		Updates:       duplicator.MakeDuplicator("update-Sender-" + userUuid),
		Deletes:       duplicator.MakeDuplicator("delete-Sender-" + userUuid),
		Notifications: duplicator.MakeDuplicator("notification-Sender-" + userUuid),
		Output:        duplicator.MakeDuplicator("output-message-" + userUuid),
		close:         make(chan interface{}),
	}
	sender.Output.RegisterInput(GlobalMessages.GetBufferedOutput(10000))
	sender.runSendDeletes()
	sender.runSendObjects()
	sender.runSendUpdates()
	sender.runSendNotifications()
	sender.Output.EnableSink()
	Senders[userUuid] = sender
	return sender
}

func GetSender(userUuid string) *Sender {
	return Senders[userUuid]
}

func (s *Sender) GetOutput() chan interface{} {
	s.lock.Acquire("new output")
	defer s.lock.Release()
	if s.activeClients == 0 {
		s.Output.DiableSink()
	}
	s.activeClients += 1
	return s.Output.GetBufferedOutput(1000)
}

func (s *Sender) CloseOutput(ch chan interface{}) {
	s.lock.Acquire("close output")
	defer s.lock.Release()
	s.Output.UnregisterOutput(ch)
	s.activeClients -= 1
	if s.activeClients == 0 {
		s.Output.EnableSink()
	}
}

func (s *Sender) RegisterUpdateInput(ch chan interface{}) {
	s.Updates.RegisterInput(ch)
}

func (s *Sender) stop() {
	for i := 0; i < 4; i++ {
		s.close <- nil
	}
	duplicator.UnlinkDouplicator(s.Output, wires.Globals)

	s.Notifications.StopDuplicator()
	s.Updates.StopDuplicator()
	s.NewObjects.StopDuplicator()
	s.Deletes.StopDuplicator()
	s.Output.StopDuplicator()

	close(s.close)
}

func (s *Sender) runSendObjects() {
	object := s.NewObjects.GetBufferedOutput(10000)
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

func (s *Sender) runSendUpdates() {
	object := s.Updates.GetBufferedOutput(1000)
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

func (s *Sender) runSendDeletes() {
	object := s.Deletes.GetBufferedOutput(1000)
	go func() {
		for {
			select {
			case newObject := <-object:
				s.Deletes.Offer(messages.BuildDeleteMessage(

					newObject.(change.Identifiable)))
			case <-s.close:
				break
			}
		}
	}()
}

func (s *Sender) runSendNotifications() {
	notifications := s.Notifications.GetBufferedOutput(1000)
	go func() {
		for {
			select {
			case newNotifications := <-notifications:
				s.Output.Offer(messages.BuildNotificationMessage(newNotifications))
			case <-s.close:
				break
			}
		}
	}()
}
