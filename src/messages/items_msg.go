package messages

import (
	"encoding/json"
	"errors"
)

/*
When adding a items
const at top
process code
JSON unmarshal
*/

const ItemAction = "item"
const BuyItemAction = "buy"
const ViewItemAction = "view"
const UseItemAction = "use"

type ItemMessage struct {
	Action string      `json:"action"`
	O      interface{} `json:"o"`
}

func (*ItemMessage) message() { return }

type ItemBuyMessage struct {
	ItemUid  string `json:"uuid"`
	ItemName string `json:"item_name"`
	Success  bool   `json:"success"`
	Err      string `json:"err,omitempty"`
}

type ItemUseMessage struct {
	ItemUuid      string      `json:"uuid"`
	UseParameters interface{} `json:"params,omitempty"`
	Result        interface{} `json:"result,omitempty"`
	Success       bool        `json:"success"`
	Err           string      `json:"err,omitempty"`
}

type ItemViewMessage struct {
	ItemUuid string      `json:"uuid"`
	Result   interface{} `json:"result,omitempty"`
	Success  bool        `json:"success"`
	Err      string      `json:"err"`
}

func BuildItemBuySuccessMessage(itemName, requestId, itemUuid string) *BaseMessage {
	return &BaseMessage{
		Action: ResponseAction,
		Msg: &ItemMessage{
			Action: BuyItemAction,
			O: &ItemBuyMessage{
				ItemName: itemName,
				Success:  true,
				ItemUid:  itemUuid,
			},
		},
		RequestID: requestId,
	}
}

func BuildItemBuyFailedMessage(itemName, requestId string, err error) *BaseMessage {
	return &BaseMessage{
		Action: ResponseAction,
		Msg: &ItemMessage{
			Action: BuyItemAction,
			O: &ItemBuyMessage{
				ItemName: itemName,
				Success:  false,
				Err:      err.Error(),
			},
		},
		RequestID: requestId,
	}
}

func BuildItemUseFailedMessage(itemUuid, requestId string, err error) *BaseMessage {
	return &BaseMessage{
		Action: ResponseAction,
		Msg: &ItemMessage{
			Action: UseItemAction,
			O: &ItemUseMessage{
				ItemUuid: itemUuid,
				Success:  false,
				Err:      err.Error(),
			},
		},
		RequestID: requestId,
	}
}

func BuildItemUseMessage(itemUUid, requestId string, result interface{}) *BaseMessage {
	return &BaseMessage{
		Action: ResponseAction,
		Msg: &ItemMessage{
			Action: UseItemAction,
			O: &ItemUseMessage{
				ItemUuid: itemUUid,
				Success:  true,
				Result:   result,
			},
		},
		RequestID: requestId,
	}
}

func BuildItemViewFailedMessage(itemUuid, requestId string, err error) *BaseMessage {
	return &BaseMessage{
		Action: ResponseAction,
		Msg: &ItemMessage{
			Action: ViewItemAction,
			O: &ItemViewMessage{
				ItemUuid: itemUuid,
				Success:  false,
				Err:      err.Error(),
			},
		},
		RequestID: requestId,
	}
}

func BuildItemViewMessage(itemUUid, requestId string, result interface{}) *BaseMessage {
	return &BaseMessage{
		Action: ResponseAction,
		Msg: &ItemMessage{
			Action: ViewItemAction,
			O: &ItemViewMessage{
				ItemUuid: itemUUid,
				Success:  true,
				Result:   result,
			},
		},
		RequestID: requestId,
	}
}

/**
custom unmarshal for json data since the action depends on what the lower level msg is
*/
func (itemMessage *ItemMessage) UnmarshalJSON(data []byte) error {
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

	if _, ok := obj["o"].(string); ok {
		return errors.New("value not there")
	}

	// see what type of action message we should use
	var o interface{}
	switch actionType {
	case BuyItemAction:
		o = &ItemBuyMessage{}
	case UseItemAction:
		o = &ItemUseMessage{}
	case ViewItemAction:
		o = &ItemViewMessage{}
	}

	str, _ := json.Marshal(obj["o"])
	err = json.Unmarshal(str, &o)
	if err != nil {
		return err
	}
	itemMessage.Action = actionType
	itemMessage.O = o
	return nil
}
