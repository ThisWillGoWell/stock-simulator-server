package messages

import "time"

const ChatAction = "chat"

type ChatMessage struct {
	Message   string    `json:"message_body"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
}

func (*ChatMessage) message() { return }

func (baseMessage *BaseMessage) IsChat() bool {
	return baseMessage.Action == "chat"
}

func BuildChatMessage(message *ChatMessage) *BaseMessage {
	return &BaseMessage{
		Action: ChatAction,
		Msg:    message,
	}
}
