package messages

const TransferAction = "transfer"

type TransferMessage struct {
	Amount int64 `json:"amount"`
	Recipient string `json:"recipient"`
}

type TransferResponse struct{
	Transfer *TransferMessage `json:"transfer"`
	Response interface{} `json:"response"`
}

func (*TransferMessage) message() { return }

func (*TransferResponse) message() { return }


func (baseMessage *BaseMessage) IsTransfer() bool {
	return baseMessage.Action == TransferAction
}


func BuildTransferResponse(response interface{}) *BaseMessage {
	return &BaseMessage{
		Action: AlertAction,
		Msg: &AlertMessage{
			Alert: response,
			Type:  "transfer_response",
		},
	}
}