package messages

import (
	"encoding/json"
	"errors"
)

const TradeAction = "trade"
const ChatAction = "chat"
const UpdateAction = "update"
const ErrorAction = "error"
const LoginAction = "login"

const ValuableUpdate = "valuable"
const PortfolioUpdate = "portfolio"
const LedgerUpdate = "ledger"

type Message interface {
	message()
}

type BaseMessage struct {
	Action string      `json:"action"`
	Msg    interface{} `json:"msg"`
}

func (baseMessage *BaseMessage) IsChat() bool {
	return baseMessage.Action == "chat"
}

func (baseMessage *BaseMessage) IsLogin() bool {
	return baseMessage.Action == LoginAction
}

func (baseMessage *BaseMessage) IsUpdate() bool {
	return baseMessage.Action == "update"
}

func (baseMessage *BaseMessage) IsTrade() bool {
	return baseMessage.Action == "trade"
}

type ErrorMessage struct {
	Err string `json:"error"`
}

func (*ErrorMessage) message() { return }

func NewErrorMessage(err string) *ErrorMessage {
	return &ErrorMessage{
		Err: err,
	}
}

type BatchMessage []BaseMessage

type LoginMessage struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (*LoginMessage) message() { return }

type TradeMessage struct {
	StockTicker string  `json:"stock_ticker"`
	ExchangeID  string  `json:"exchange_id"`
	Amount      float64 `json:"amount"`
}

func (*TradeMessage) message() { return }

type TradeResponse struct {
	Trade    *TradeMessage `json:"trade"`
	Response interface{}   `json:"response"`
}

func (*TradeResponse) message() { return }

func BuildPurchaseResponse(message *TradeMessage, response interface{}) *TradeResponse {
	return &TradeResponse{
		Trade:    message,
		Response: response,
	}
}

type ChatMessage struct {
	Message   string `json:"message_body"`
	Author    string `json:"author"`
	Timestamp int64  `json:"timestamp"`
}

func (*ChatMessage) message() { return }

type UpdateMessage struct {
}

func (*UpdateMessage) message() { return }

func BuildUpdateMessage(obj interface{}) *BaseMessage {
	return &BaseMessage{
		Action: UpdateAction,
		Msg:    &obj,
	}
}

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

	if _, ok := obj["value"].(string); ok {
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
	}

	str, _ := json.Marshal(obj["value"])
	err = json.Unmarshal(str, &message)
	if err != nil {
		return err
	}
	baseMessage.Action = actionType
	baseMessage.Msg = message

	return nil
}
