package messages

const UpdateAction = "update"

type UpdateMessage struct {
}

func (*UpdateMessage) message() { return }

func BuildUpdateMessage(obj interface{}) *BaseMessage {
	return &BaseMessage{
		Action: UpdateAction,
		Msg:    &obj,
	}
}

func (baseMessage *BaseMessage) IsUpdate() bool {
	return baseMessage.Action == "update"
}
