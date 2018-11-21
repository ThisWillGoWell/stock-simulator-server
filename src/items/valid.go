package items

import (
	"encoding/json"

	"github.com/stock-simulator-server/src/log"
)

func LoadItemConfig(data []byte) {
	validList := make([]*validItemConfiguration, 0)
	err := json.Unmarshal(data, &validList)
	if err != nil {
		log.Log.Error("error reading items config ", err)
	}
	for _, ele := range validList {
		validItems[ele.ConfigId] = ele
	}
}

var validItems = make(map[string]*validItemConfiguration)

// also change vic2 and the unmarshal
type validItemConfiguration struct {
	ConfigId      string    `json:"config_id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Cost          int64     `json:"cost"`
	RequiredLevel int64     `json:"required_level"`
	Prams         InnerItem `json:"params"`
}
type vic2 struct {
	ConfigId      string      `json:"config_id"`
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	Cost          int64       `json:"cost"`
	RequiredLevel int64       `json:"required_level"`
	Prams         interface{} `json:"params"`
}

func (config *validItemConfiguration) UnmarshalJSON(data []byte) error {
	vic2 := new(vic2)

	err := json.Unmarshal(data, vic2)
	config.ConfigId = vic2.ConfigId
	config.Name = vic2.Name
	config.Type = vic2.Type
	config.Cost = vic2.Cost
	config.RequiredLevel = vic2.RequiredLevel

	if err != nil {
		return err
	}
	pStr, _ := json.Marshal(vic2.Prams)

	switch config.Type {
	case TradeItemType:
		itemPointer := &TradeEffectItem{}
		json.Unmarshal(pStr, itemPointer)
		config.Prams = itemPointer

	}

	return nil
}
