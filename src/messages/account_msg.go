package messages

const LoginAction = "login"
const NewAccountAction = "new_account"

type LoginMessage struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (*LoginMessage) message() { return }

func (baseMessage *BaseMessage) IsAccountCreate() bool {
	return baseMessage.Action == NewAccountAction
}

type NewAccountMessage struct {
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func (*NewAccountMessage) message() { return }

func (baseMessage *BaseMessage) IsLogin() bool {
	return baseMessage.Action == LoginAction
}

type AccountResponseMessage struct {
	Success bool   `json:"success"`
	Uuid    string `json:"uuid"`
	Err     string `json:"err"`
}

func (*AccountResponseMessage) message() { return }

func SuccessLogin(userGuid string) *BaseMessage {
	return &BaseMessage{
		Action: LoginAction,
		Msg: &AccountResponseMessage{
			Success: true,
			Uuid:    userGuid,
			Err:     "",
		},
	}
}

func FailedLogin(err error) *BaseMessage {
	return &BaseMessage{
		Action: LoginAction,
		Msg: &AccountResponseMessage{
			Success: false,
			Uuid:    "",
			Err:     err.Error(),
		},
	}
}
