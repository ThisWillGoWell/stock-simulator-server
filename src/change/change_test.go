package change

import (
	"testing"
)

type changeTest struct {
	Name string
	Test int `change:"-"`
}

func (change *changeTest) GetId() string {
	return change.Name
}

func (change *changeTest) GetType() string {
	return "test"
}

func TestChange(t *testing.T) {
	StartDetectChanges()

	PublicSubscribeChange.EnableCopyMode()
	messagesReceived := 0
	done := make(chan interface{})
	go func() {
		changes := PublicSubscribeChange.GetBufferedOutput(2)
		for c := range changes {
			t.Log(c)
			messagesReceived += 1
			if messagesReceived == 3 {
				done <- nil
			}

		}
	}()

	v := &changeTest{
		Name: "this is my name",
		Test: 1,
	}
	SubscribeUpdateInputs.Offer(v)
	v.Test = 2
	SubscribeUpdateInputs.Offer(v)
	v.Test = 3
	SubscribeUpdateInputs.Offer(v)

	<-done

}
