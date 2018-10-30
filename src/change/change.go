package change

import (
	"fmt"
	"reflect"

	"github.com/stock-simulator-server/src/log"

	"github.com/stock-simulator-server/src/wires"

	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/lock"
)

/*
	So this is some of the most bs code I have written in awhile
	So how can we detect the change in a pointer's value in any struct?
	welp this works
	The idea is to use tags of a struct to watch changes on a value using reflect
	and then using the Identifiable Interface we can tag the change


	type Identifiable interface {
		Id() string
		GetType() string
	}

	type changeTest struct {
			Name	string
			Test int `change:"-"`
	}

	func (change *changeTest)Id() string{
		return change.Name
	}

	func (change *changeTest)GetType() string{
		return "test"
	}


	Here the change detect is registered and will be tagged with the Id Field and its value

	&changeTest{
		Id: "this is my id",
		Test: 1,
	}

	[{
		"id": "this is my name
		"type": "test"
		"name": "Test",
		"value": 2
	}]

*/

const (
	changeTag = "change"
)

var (
	// subscribeables is something that can be subscribed to
	subscribeables        = make(map[string]*SubscribeUpdate)
	subscribeablesLock    = lock.NewLock("subscribeables")
	PublicSubscribeChange = duplicator.MakeDuplicator("subscribe-outputs")
)

func RegisterPublicChangeDetect(o Identifiable) error {
	log.Log.Trace("registering public change detect: ", o.GetType(), o.GetId())
	output := make(chan interface{})
	PublicSubscribeChange.RegisterInput(output)
	return registerChangeDetect(o, output)
}

func RegisterPrivateChangeDetect(o Identifiable, update chan interface{}) error {
	log.Log.Trace("registering private change detect", o.GetType(), o.GetId())
	return registerChangeDetect(o, update)
}

func UnregisterChangeDetect(o Identifiable) {
	subscribeablesLock.Acquire("unregister-change")
	defer subscribeablesLock.Release()
	if _, ok := subscribeables[o.GetType()+o.GetId()]; !ok {
		log.Alerts.Fatal("Panic in Change Detect, cant unregister since does not exists change", o.GetId(), o.GetId())
		log.Log.Fatal("Panic in Change Detect, cant unregister since does not exists change", o.GetId(), o.GetId())
		panic("cant unregister change detect that does not exists" + o.GetType() + o.GetId())
	}
	delete(subscribeables, o.GetType()+o.GetId())
}

func registerChangeDetect(o Identifiable, outputChan chan interface{}) error {
	//get the include tags
	subscribeablesLock.Acquire("register-change")
	defer subscribeablesLock.Release()
	if o.GetId() == "291" {
		fmt.Println("")
	}

	if _, ok := subscribeables[o.GetType()+o.GetId()]; ok {
		log.Alerts.Fatal("Panic in Change Detect, cant register since already exists", o.GetId(), o.GetId())
		log.Log.Fatal("Panic in Change Detect, cant register since already existse", o.GetId(), o.GetId())
		panic("change detect already registered, check the code" + o.GetType() + o.GetId())
	}

	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = reflect.ValueOf(o).Elem().Type()
	}
	//t := reflect.ValueOf(o).Type()
	newChangeDetect := &SubscribeUpdate{
		Type:          o.GetType(),
		Id:            o.GetId(),
		changeDetects: make(map[string]*ChangeField),
		TargetOutput:  outputChan,
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		_, exists := field.Tag.Lookup(changeTag)
		if exists {
			jsonFieldName, exists := field.Tag.Lookup("json")
			if !exists {
				jsonFieldName = field.Name
			}

			changeField := &ChangeField{
				Value: nil,
				Field: jsonFieldName,
			}
			newChangeDetect.changeDetects[field.Name] = changeField
		}
	}

	subscribeables[o.GetType()+o.GetId()] = newChangeDetect
	return nil
}

func getValue(o interface{}, name string) interface{} {
	var r reflect.Value
	//get the type of the interface provided
	t := reflect.TypeOf(o)
	// if pointer, dereference it and then get the field
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		r = reflect.ValueOf(o).Elem().FieldByName(name)
	} else {
		r = reflect.ValueOf(o).FieldByName(name)
	}

	//is the value of that field a pointer?
	if r.Kind() == reflect.Ptr || r.Kind() == reflect.Interface {
		return reflect.Indirect(r)
	}
	return r.Interface()
}

func StartDetectChanges() {
	//SubscribeUpdateInputs.EnableCopyMode()
	//SubscribeUpdateInputs.EnableDebug()
	//SubscribeUpdateOutput.EnableDebug()
	allUpdates := duplicator.MakeDuplicator("all-updates")
	allUpdates.RegisterInput(wires.ItemsUpdate.GetBufferedOutput(10000))
	allUpdates.RegisterInput(wires.StocksUpdate.GetBufferedOutput(10000))
	allUpdates.RegisterInput(wires.PortfolioUpdate.GetBufferedOutput(10000))
	allUpdates.RegisterInput(wires.LedgerUpdate.GetBufferedOutput(10000))
	allUpdates.RegisterInput(wires.UsersUpdate.GetBufferedOutput(10000))
	allUpdates.RegisterInput(wires.BookUpdate.GetBufferedOutput(10000))
	PublicSubscribeChange.EnableCopyMode()
	//allUpdates.RegisterInput(wires.NotificationUpdate.GetBufferedOutput(10000))
	subscribeUpdateChannel := allUpdates.GetBufferedOutput(10000)
	go func() {
		for updateObj := range subscribeUpdateChannel {
			update, ok := updateObj.(Identifiable)
			if !ok {
				continue
			}
			subscribeablesLock.Acquire("detect change")

			changeDetect, exists := subscribeables[update.GetType()+update.GetId()]
			if !exists {
				continue
			}

			changedFields := make([]*ChangeField, 0)
			changed := false
			for filedName, fieldChange := range changeDetect.changeDetects {
				currentValue := getValue(update, filedName)

				if !reflect.DeepEqual(currentValue, fieldChange.Value) {
					changed = true
					fieldChange.Value = currentValue
					changedFields = append(changedFields, &ChangeField{fieldChange.Field, fieldChange.Value})
				}
			}
			if changed {
				changeDetect.TargetOutput <- &ChangeNotify{
					Type:    changeDetect.Type,
					Id:      changeDetect.Id,
					Changes: changedFields,
				}
			}
			subscribeablesLock.Release()
		}
	}()
}

type SubscribeUpdate struct {
	Type          string
	Id            string
	changeDetects map[string]*ChangeField
	TargetOutput  chan interface{} `json:"-"`
}

type ChangeField struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type ChangeNotify struct {
	Type    string         `json:"type"`
	Id      string         `json:"uuid"`
	Changes []*ChangeField `json:"changes"`
	Object  Identifiable   `json:"-"`
}

func (cn *ChangeNotify) GetId() string {
	return cn.Id
}

func (cn *ChangeNotify) GetType() string {
	return "Change-Notify"
}

func GetCurrentValues() []*ChangeNotify {
	subscribeablesLock.Acquire("current values")
	defer subscribeablesLock.Release()
	values := make([]*ChangeNotify, 0)
	for _, value := range subscribeables {
		currentVals := make([]*ChangeField, 0)
		for _, val := range value.changeDetects {
			currentVals = append(currentVals, val)
		}

		newVal := &ChangeNotify{
			Type:    value.Type,
			Id:      value.Id,
			Changes: currentVals,
		}
		values = append(values, newVal)
	}
	return values
}

func Test() {

}
