package messages

import "github.com/stock-simulator-server/src/change"

const ObjectAction = "object"

type ObjectMessage struct {
	Type string `json:"type"`
	Id string `json:"uuid"`
	Value change.Identifiable `json:"object"`
}
func (*ObjectMessage) message() { return }

func NewObjectMessage(identifiable change.Identifiable)*BaseMessage{
	return &BaseMessage{
		Action:ObjectAction,
		Msg:ObjectMessage{
			Type: identifiable.GetType(),
			Id: identifiable.GetId(),
			Value: identifiable,
		},
	}
}