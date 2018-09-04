package messages

import (
	"encoding/json"
	"errors"
	"time"
)

const QueryAction = "query"

func (baseMessage *BaseMessage) IsQuery() bool {
	return baseMessage.Action == QueryAction
}

type QueryMessage struct {
	QueryUUID string    `json:"uuid"`
	QueryField string   `json:"field"`
	NumberPoints int    `json:"num_points"`
	Length Duration 	`json:"length"`
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type QueryResponse struct{
	Success bool `json:"success"`
	Error string `json:"error"`
	Message *QueryMessage `json:"message"`
	Points [][]interface{} `json:"points"`

}

func (*QueryMessage) message() { return }
