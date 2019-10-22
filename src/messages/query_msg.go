package messages

import "github.com/ThisWillGoWell/stock-simulator-server/src/utils"

const QueryAction = "query"

func (baseMessage *BaseMessage) IsQuery() bool {
	return baseMessage.Action == QueryAction
}

type QueryMessage struct {
	QueryUUID     string         `json:"uuid"`
	QueryField    string         `json:"field"`
	NumberPoints  int            `json:"num_points"`
	Length        utils.Duration `json:"length"`
	ForceUpdate   bool           `json:"force_update"`
	CacheDuration utils.Duration `json:"cache_duration"`
}

type QueryResponse struct {
	Success bool            `json:"success"`
	Error   string          `json:"error"`
	Message *QueryMessage   `json:"message"`
	Points  [][]interface{} `json:"points"`
}

func (*QueryMessage) message() { return }
