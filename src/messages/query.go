package messages

import "time"

const QueryAction = "query"

func (baseMessage *BaseMessage) IsQuery() bool {
	return baseMessage.Action == QueryAction
}

type QueryMessage struct {
	QueryUUID string    `json:"uuid"`
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
}

func (*QueryMessage) message() { return }
