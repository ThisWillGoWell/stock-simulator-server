package items

import (
	"encoding/json"

	"github.com/stock-simulator-server/src/log"
)

func LoadItemConfig(data []byte) {
	err := json.Unmarshal(data, &validItems)
	if err != nil {
		log.Log.Error("error reading items config ", err)
	}
}

var validItems = make(map[string]validItemConfiguration)

type validItemConfiguration struct {
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Cost          int64     `json:"cost"`
	RequiredLevel int64     `json:"required_level"`
	Prams         InnerItem `json:"params"`
}

func (config *validItemConfiguration) UnmarshalJSON(data []byte) error {

	err := json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	pStr, _ := json.Marshal(config.Prams)

	switch config.Type {
	case TradeItemType:
		config.Prams = new(TradeEffectItem)

	}
	json.Unmarshal(pStr, config.Prams)

	return nil
}
