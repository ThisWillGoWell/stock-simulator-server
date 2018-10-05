package messages

const ConnectAction = "connect"

type ConnectMessage struct {
	SessionToken string `json:"token"`
}

func (*ConnectMessage) message() { return }

func (baseMessage *BaseMessage) IsConnect() bool {
	return baseMessage.Action == ConnectAction
}

type AccountResponseMessage struct {
	Success      bool                   `json:"success"`
	Config       map[string]interface{} `json:"config"`
	SessionToken string                 `json:"token,omitempty"`
	Uuid         string                 `json:"uuid,omitempty"`
	Err          string                 `json:"err,omitempty"`
	Init         map[string]interface{} `json:"init,omitempty"`
}

func (*AccountResponseMessage) message() { return }

func SuccessConnect(userGuid, token string, config map[string]interface{}) *BaseMessage {
	return &BaseMessage{
		Action: ConnectAction,
		Msg: &AccountResponseMessage{
			Success:      true,
			SessionToken: token,
			Uuid:         userGuid,
			Err:          "",
			Config:       config,
		},
	}
}

func FailedConnect(err error) *BaseMessage {
	return &BaseMessage{
		Action: ConnectAction,
		Msg: &AccountResponseMessage{
			Success: false,
			Uuid:    "",
			Err:     err.Error(),
		},
	}
}
