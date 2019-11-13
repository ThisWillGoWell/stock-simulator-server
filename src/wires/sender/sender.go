package sender

import (
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id/change"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/messages"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/duplicator"
)

var senders = make(map[string]*Sender)

var GlobalMessages = duplicator.MakeDuplicator("global-messages")

func RunGlobalSender() {
	go func() {
		globalObjects := duplicator.MakeDuplicator("global-new-objects")
		globalObjects.RegisterInput(wires.UsersNewObject.GetBufferedOutput(10000))
		globalObjects.RegisterInput(wires.StocksNewObject.GetBufferedOutput(10000))
		globalObjects.RegisterInput(wires.PortfolioNewObject.GetBufferedOutput(10000))
		globalObjects.RegisterInput(wires.LedgerNewObject.GetBufferedOutput(10000))
		globalObjects.RegisterInput(wires.BookNewObject.GetBufferedOutput(10000))
		globalObjects.RegisterInput(wires.RecordsNewObject.GetBufferedOutput(10000))
		globalObjects.RegisterInput(wires.EffectsNewObject.GetBufferedOutput(10000))
		out := globalObjects.GetBufferedOutput(100000)
		for ele := range out {
			GlobalMessages.Offer(messages.NewObjectMessage(ele.(id.Identifiable)))
		}
	}()

	go func() {
		out := change.PublicSubscribeChange.GetBufferedOutput(100000)
		for ele := range out {
			GlobalMessages.Offer(messages.BuildUpdateMessage(ele.(id.Identifiable)))
		}
	}()

	go func() {
		globalDeletes := duplicator.MakeDuplicator("global-deletes")
		globalDeletes.RegisterInput(wires.EffectsDelete.GetBufferedOutput(10000))
		for ele := range globalDeletes.GetBufferedOutput(10000) {
			GlobalMessages.Offer(messages.BuildDeleteMessage(ele.(id.Identifiable)))
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
	userId        string
	portfolioId   string
}

func NewSender(userUuid, portfolioUuid string) *Sender {
	sender := &Sender{
		lock:          lock.NewLock("client-user-Sender-" + userUuid),
		activeClients: 0,
		NewObjects:    duplicator.MakeDuplicator("objects-Sender-" + userUuid),
		Updates:       duplicator.MakeDuplicator("update-Sender-" + userUuid),
		Deletes:       duplicator.MakeDuplicator("delete-Sender-" + userUuid),
		Notifications: duplicator.MakeDuplicator("notification-Sender-" + userUuid),
		Output:        duplicator.MakeDuplicator("output-message-" + userUuid),
		close:         make(chan interface{}),
		userId:        userUuid,
		portfolioId:   portfolioUuid,
	}
	sender.Output.RegisterInput(GlobalMessages.GetBufferedOutput(100))
	sender.runSendDeletes()
	sender.runSendObjects()
	sender.runSendUpdates()
	sender.runSendNotifications()
	senders[userUuid] = sender
	senders[portfolioUuid] = sender
	return sender
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

func (s *Sender) Stop() {
	s.Notifications.StopDuplicator()
	s.Updates.StopDuplicator()
	s.NewObjects.StopDuplicator()
	s.Deletes.StopDuplicator()
	s.Output.StopDuplicator()
	delete(senders, s.portfolioId)
	delete(senders, s.userId)
	duplicator.UnlinkDuplicator(s.Output, GlobalMessages)
	close(s.close)
}

func (s *Sender) runSendObjects() {
	object := s.NewObjects.GetBufferedOutput(10000)
	go func() {
		for {
			select {
			case newObject := <-object:
				s.Output.Offer(messages.NewObjectMessage(newObject.(id.Identifiable)))
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
				s.Output.Offer(messages.BuildUpdateMessage(newObject.(id.Identifiable)))
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
				s.Output.Offer(messages.BuildDeleteMessage(
					newObject.(id.Identifiable)))
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

func SendNewObject(id string, newObject id.Identifiable) {
	if _, exists := senders[id]; !exists {
		fmt.Println("cant find sender id during new: " + id)
		return
	}
	senders[id].NewObjects.Offer(newObject)
}

func SendDeleteObject(id string, deleteObject id.Identifiable) {
	if _, exists := senders[id]; !exists {
		fmt.Println("cant find sender id during delete: " + id)
		return
	}
	senders[id].Deletes.Offer(deleteObject)
}

func RegisterChangeUpdate(id string, changeChannel chan interface{}) error {
	if _, exists := senders[id]; !exists {
		return fmt.Errorf("cant find sender during add change update: " + id)

	}
	senders[id].Updates.RegisterInput(changeChannel)
	return nil
}

func SendChangeUpdate(id string, changeNotify *change.ChangeNotify) {
	if _, exists := senders[id]; !exists {
		fmt.Println("cant find sender during add change update: " + id)
		return
	}
	senders[id].Updates.Offer(changeNotify)
}
