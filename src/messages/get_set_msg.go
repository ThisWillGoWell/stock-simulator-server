package messages

const SetAction = "set"

type SetMessage struct {
	Set string `json:"set"`
	Value interface{} `json:"value"`
}

type SetResponse struct {
	Success bool `json:"success"`
	Err string `json:"error,omitempty"`
}
func (*SetMessage) message() { return }

func BuildSuccessSet() *SetResponse{
	return &SetResponse{
		Success: true,
	}
}

func BuildFailedSet(err error) *SetResponse{
	return  &SetResponse{
		Success: false,
		Err: err.Error(),
	}
}
