package messages

const ChatAction = "chat"

type ChatMessage struct {
	Message   string `json:"message_body"`
	Author    string `json:"author"`
	Timestamp int64  `json:"timestamp"`
}
func (*ChatMessage) message() { return }

func (baseMessage *BaseMessage) IsChat() bool {
	return baseMessage.Action == "chat"
}
