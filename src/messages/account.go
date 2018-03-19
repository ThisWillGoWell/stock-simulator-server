package messages

const LoginAction = "login"
const NewAccountAction="new_account"

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