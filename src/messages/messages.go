package messages

import (
	"encoding/json"
	"errors"
)

type Message interface {
	message()
}

type BaseMessage struct {
	Action string      `json:"action"`
	Msg    interface{} `json:"msg"`
}

type BatchMessage []BaseMessage

/**
custom unmarshal for json data since the action depends on what the lower level msg is
*/
func (baseMessage *BaseMessage) UnmarshalJSON(data []byte) error {
	//start with a generic string -> interface map
	var obj map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	// make sure that generic one contains the required keys
	actionType := ""
	if t, ok := obj["action"].(string); ok {
		actionType = t
	} else {
		return errors.New("action not there")
	}

	if _, ok := obj["msg"].(string); ok {
		return errors.New("value not there")
	}

	// see what type of action message we should use
	// update is not here since it should never have to be Unmarshal
	var message Message
	switch actionType {
	case ChatAction:
		message = &ChatMessage{}
	case TradeAction:
		message = &TradeMessage{}
	case LoginAction:
		message = &LoginMessage{}
	case UpdateAction:
		message = &UpdateMessage{}
	case NewAccountAction:
		message = &NewAccountMessage{}
	case QueryAction:
		message = &QueryMessage{}
	case TransferAction:
		message = &TransferMessage{}

	}

	str, _ := json.Marshal(obj["msg"])
	err = json.Unmarshal(str, &message)
	if err != nil {
		return err
	}
	baseMessage.Action = actionType
	baseMessage.Msg = message

	return nil
}
