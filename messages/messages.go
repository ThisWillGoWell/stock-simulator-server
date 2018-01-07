package messages

import (
	"encoding/json"
	"stock-server/utils"
	"stock-server/client"
	"errors"
)

const TradeAction  = "trade"
const ChatAction  = "chat"
const UpdateAction = "update"

const StockUpdate = "stock"
const PortfolioUpdate = "portfolio"

type MessageProcessor struct {
	inputchannel chan string
	ChatMessages utils.ChannelDuplicator
	TradeCommands utils.ChannelDuplicator
}

func RunMessageProcessor() {

}
type Message interface {
	message()
}



type BaseMessage struct{
	Action string `json:"action"`
	Value  Message `json:"value"`
	client client.Client
}
func (msg *BaseMessage) IsChat()bool{
	return msg.Action == "chat"
}

func (msg *BaseMessage) IsUpdate()bool{
	return msg.Action == "update"
}

func (msg *BaseMessage) IsTrade()bool{
	return msg.Action == "trade"
}



type TradeMessage struct {
	StockTicker string `json:"stock_ticker"`
	Amount float64 `json:"amount"`
}
func (*TradeMessage) message() {return}

type ChatMessage struct {
	Message string `json:"message_body"`
	Author string `json:"author"`
	Timestamp int64 `json:"timestamp"`
}
func (*ChatMessage) message() {return}


type UpdateMessage struct {
	UpdateType string `json:"type"`
	Object interface{} `json:"object"`
}
func (*UpdateMessage) message() {return}

func NewUpdateMessage(t string, obj interface{})*BaseMessage{
	updateMsg :=  UpdateMessage{
		UpdateType: t,
		Object:obj,
	}
	return &BaseMessage{
		Action: UpdateAction,
		Value: &updateMsg,
	}
}

func (baseMessage *BaseMessage)UnmarshalJSON(data [] byte) error{
	//start with a generic string -> interface map
	var obj map[string]interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil{
		return err
	}
	// make sure that generic one contains the required keys
	actionType := ""
	if t, ok := obj["action"].(string); ok {
		actionType = t
	}else{
		return errors.New("action not there")
	}

	if _, ok := obj["value"].(string); ok{
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
	}

	str,_ := json.Marshal(obj["value"])
	err = json.Unmarshal(str, &message)
	if err != nil {
		return err
	}
	*baseMessage = BaseMessage{
		Action: actionType,
		Value:  message,
	}
	return nil
}

func main(){

}