package change

import (
	"encoding/json"
	"testing"

	"github.com/stock-simulator-server/src/deepcopy"

	"github.com/stock-simulator-server/src/wires"
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

func waitForChanges(amount int, t *testing.T) chan interface{} {
	messagesReceived := 0
	done := make(chan interface{})
	go func() {
		PublicSubscribeChange.EnableDebug()
		changes := PublicSubscribeChange.GetBufferedOutput(2)
		for c := range changes {
			change, _ := json.Marshal(c)
			t.Log(string(change))
			messagesReceived += 1
			if messagesReceived == amount {
				done <- nil
			}

		}
	}()
	return done
}

func TestChange(t *testing.T) {

	PublicSubscribeChange.EnableCopyMode()
	done := waitForChanges(3, t)
	v := &changeTest{
		Name: "this is my name",
		Test: 1,
	}
	wires.ConnectWires()
	RegisterPublicChangeDetect(v)

	wires.ItemsUpdate.Offer(v)
	v.Test = 2
	wires.ItemsUpdate.Offer(v)
	v.Test = 3
	wires.ItemsUpdate.Offer(v)

	<-done
}

type ArrayChange struct {
	ID  string
	Arr []string `change:"-"`
}

func (ac *ArrayChange) GetId() string {
	return ac.ID
}
func (*ArrayChange) GetType() string {
	return "array-change"
}

func TestArrayChange(t *testing.T) {
	StartDetectChanges()
	done := waitForChanges(2, t)
	v := &ArrayChange{
		"1",
		[]string{"hello", "world"},
	}
	wires.ConnectWires()
	RegisterPublicChangeDetect(v)
	wires.ItemsUpdate.Offer(v)
	v.Arr = []string{"world", "hello"}
	wires.ItemsUpdate.Offer(v)
	<-done
}

type Inner struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}
type ArrayStruct struct {
	ID  string  `json:"id"`
	Arr []Inner `json:"array" change:"-"`
}

func (ac *ArrayStruct) GetId() string {
	return ac.ID
}
func (*ArrayStruct) GetType() string {
	return "array-struct"
}

func TestArrayStructChange(t *testing.T) {
	StartDetectChanges()
	done := waitForChanges(2, t)
	v := &ArrayStruct{
		"1",
		[]Inner{{"name1", 1}, {"name2", 2}},
	}
	newV := deepcopy.Copy(v)
	_, ok := newV.(interface{})
	for ok {
		break
	}
	wires.ConnectWires()
	RegisterPublicChangeDetect(v)
	wires.ItemsUpdate.Offer(v)
	v.Arr = v.Arr[:1]
	wires.ItemsUpdate.Offer(v)
	<-done
}
