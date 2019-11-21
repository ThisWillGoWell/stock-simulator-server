package messages

import (
	"github.com/ThisWillGoWell/stock-simulator-server/src/id"
)

const ObjectAction = "object"

type ObjectMessage struct {
	Type  string          `json:"type"`
	Id    string          `json:"uuid"`
	Value id.Identifiable `json:"object"`
}

func (*ObjectMessage) message() { return }

func NewObjectMessage(identifiable id.Identifiable) *BaseMessage {
	return &BaseMessage{
		Action: ObjectAction,
		Msg: ObjectMessage{
			Type:  identifiable.GetType(),
			Id:    identifiable.GetId(),
			Value: identifiable,
		},
	}
}
