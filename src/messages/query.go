package messages

const QueryAction = "query"

func (baseMessage *BaseMessage) IsQuery() bool {
	return baseMessage.Action == QueryAction
}

type QueryMessage struct {
	QueryUUID string `query_id`
	Type      string `json:"historical"`
}

func (*QueryMessage) message() { return }

func NewQueryResponse([]) *QueryMessage {
	return &QueryMessage{
		Type:      "query_response",
		Alert:     err,
		Timestamp: 0,
	}
}
