package messages

import "github.com/stock-simulator-server/src/histroy"

const HistoricalAction = "historical"

type HistoricalMessage struct {
	Uuid string                     `json:"uuid"`
	Ts   []*histroy.TimeSeriesEntry `json:"ts"`
}

func (*HistoricalMessage) message() { return }

func BuildHistoericalMessage(object histroy.TimeSeriesObject) *BaseMessage {
	return &BaseMessage{
		Action: HistoricalAction,
		Msg: HistoricalMessage{
			Uuid: object.Uuid,
			Ts:   object.Data,
		},
	}
}

func (baseMessage *BaseMessage) IsHistorical() bool {
	return baseMessage.Action == HistoricalAction
}
