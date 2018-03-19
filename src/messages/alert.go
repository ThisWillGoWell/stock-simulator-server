package messages
const AlertAction = "alert"

func (baseMessage *BaseMessage) IsAlert() bool {
	return baseMessage.Action == AlertAction
}

type AlertMessage struct {
	Alert   string `json:"alert"`
	Type    string `json:"type"`
	Timestamp int64  `json:"timestamp"`
}
func (*AlertMessage) message() { return }


func NewErrorMessage(err string) *AlertMessage {
	return &AlertMessage{
		Type: "error",
		Alert: err,
		Timestamp: 0,
	}
}

