package messages

const AlertAction = "alert"

func (baseMessage *BaseMessage) IsAlert() bool {
	return baseMessage.Action == AlertAction
}

type AlertMessage struct {
	Alert     interface{} `json:"alert"`
	Type      string      `json:"type"`
}

func (*AlertMessage) message() { return }

func NewErrorMessage(err string) *AlertMessage {
	return &AlertMessage{
		Type:      "error",
		Alert:     err,
	}
}
