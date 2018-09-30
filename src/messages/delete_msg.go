package messages

import "github.com/stock-simulator-server/src/change"

const DeleteAction = "delete"

func (baseMessage *BaseMessage) isDelete() bool {
	return baseMessage.Action == DeleteAction
}

type DeleteMsg struct {
	Uuid string `json:"uuid"`
	Type string `json:"type"`
}

func (*DeleteMsg) message() { return }

func BuildDeleteMessage(o change.Identifiable) *BaseMessage {
	return &BaseMessage{
		Action: DeleteAction,
		Msg: &DeleteMsg{
			Uuid: o.GetId(),
			Type: o.GetType(),
		},
	}
}
