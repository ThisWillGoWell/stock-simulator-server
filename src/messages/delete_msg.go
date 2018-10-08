package messages

import "github.com/stock-simulator-server/src/change"

const DeleteAction = "delete"

func (baseMessage *BaseMessage) isDelete() bool {
	return baseMessage.Action == DeleteAction
}

type DeleteMessage struct {
	Uuid string `json:"uuid"`
	Type string `json:"type"`
}

type DeleteResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"err"`
}

func (*DeleteMessage) message() { return }

func BuildDeleteMessage(o change.Identifiable) *BaseMessage {
	return &BaseMessage{
		Action: DeleteAction,
		Msg: &DeleteMessage{
			Uuid: o.GetId(),
			Type: o.GetType(),
		},
	}
}
func BuildDeleteResponseMsg(requestId string, err error) *BaseMessage {
	if err == nil {
		return &BaseMessage{
			Action: ResponseAction,
			Msg: DeleteResponse{
				true,
				"",
			},
			RequestID: requestId,
		}
	}
	return &BaseMessage{
		Action: ResponseAction,
		Msg: DeleteResponse{
			false,
			err.Error(),
		},
		RequestID: requestId,
	}
}
