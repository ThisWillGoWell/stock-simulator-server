package messages

import (
	"encoding/json"
	"errors"
)

type Message interface {
	message()
}

type BaseMessage struct {
	Action    string      `json:"action"`
	Msg       interface{} `json:"msg"`
	RequestID string      `json:"request_id,omitempty"`
}

const ResponseAction = "response"

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
	requestId := ""
	if t, ok := obj["action"].(string); ok {
		actionType = t
	} else {
		return errors.New("action not there")
	}

	if t, ok := obj["request_id"].(string); ok {
		requestId = t
	} else {

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
	case RenewAction:
		message = &RenewMessage{}
	case SetAction:
		message = &SetMessage{}
	}

	str, _ := json.Marshal(obj["msg"])
	err = json.Unmarshal(str, &message)
	if err != nil {
		return err
	}
	baseMessage.Action = actionType
	baseMessage.Msg = message
	baseMessage.RequestID = requestId

	return nil
}

func BuildResponseMsg(response interface{}, requestID string) *BaseMessage {
	return &BaseMessage{
		Action:    ResponseAction,
		Msg:       response,
		RequestID: requestID,
	}
}
