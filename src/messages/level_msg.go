package messages

const LevelUpAction = "level_up"

type LevelUpMessage struct {
}

func (*LevelUpMessage) message() { return }

type LevelUpResponse struct {
	Success bool   `json:"success"`
	Err     string `json:"err,omitempty"`
}

func (*LevelUpResponse) message() { return }

func BuildLevelUpResponse(requestId string, err error) *BaseMessage {
	success := err == nil
	errMsg := ""
	if !success {
		errMsg = err.Error()
	}
	return &BaseMessage{
		Action:    ResponseAction,
		RequestID: requestId,
		Msg: &LevelUpResponse{
			Success: success,
			Err:     errMsg,
		},
	}
}
