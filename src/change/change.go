package change

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ThisWillGoWell/stock-simulator-server/src/id"

	"github.com/ThisWillGoWell/stock-simulator-server/src/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"github.com/ThisWillGoWell/stock-simulator-server/src/duplicator"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
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

func RegisterPublicChangeDetect(o id.Identifiable) error {
	log.Log.Trace("registering public change detect: ", o.GetType(), o.GetId())
	output := make(chan interface{})
	PublicSubscribeChange.RegisterInput(output)
	return registerChangeDetect(o, output)
}

func RegisterPrivateChangeDetect(o id.Identifiable, update chan interface{}) error {
	log.Log.Trace("registering private change detect", o.GetType(), o.GetId())
	return registerChangeDetect(o, update)
}

func UnregisterChangeDetect(o id.Identifiable) {
	subscribeablesLock.Acquire("unregister-change")
	defer subscribeablesLock.Release()
	if _, ok := subscribeables[o.GetType()+o.GetId()]; !ok {
		log.Log.Errorf("err in Change Detect, cant unrested since does not exists change", o.GetId(), o.GetId())
		return
	}
	delete(subscribeables, o.GetType()+o.GetId())
}

func registerChangeDetect(o id.Identifiable, outputChan chan interface{}) error {
	//get the include tags
	subscribeablesLock.Acquire("register-change")
	defer subscribeablesLock.Release()

	if _, ok := subscribeables[o.GetType()+o.GetId()]; ok {
		return fmt.Errorf("change detect already registered %s-%s", o.GetType(), o.GetId())
	}

	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = reflect.ValueOf(o).Elem().Type()
	}
	//t := reflect.ValueOf(o).Type()
	newChangeDetect := &SubscribeUpdate{
		Type:          o.GetType(),
		Id:            o.GetId(),
		changeDetects: getAllFields(o),
		TargetOutput:  outputChan,
	}

	subscribeables[o.GetType()+o.GetId()] = newChangeDetect
	return nil
}

/*
so want to be able to recursively build objects that are represented as
n-dimensional (though currently only do 2) in the system but 1-dimensional external

*/

func getAllFields(o interface{}) map[string]interface{} {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = reflect.ValueOf(o).Elem().Type()
		o = reflect.ValueOf(o).Elem().Interface()
	}

	changes := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value, exists := field.Tag.Lookup(changeTag)
		if exists {
			if value == "inner" {

				changes[field.Name] = getAllFields(reflect.ValueOf(o).FieldByName(field.Name).Interface())
			} else {
				jsonFieldName, exists := field.Tag.Lookup("json")

				if !exists {
					jsonFieldName = field.Name
				} else {
					commaIndex := strings.Index(jsonFieldName, ",")
					if commaIndex != -1 {
						jsonFieldName = jsonFieldName[:commaIndex]
					}
				}

				changeField := &ChangeField{
					Value: getValue(o, field.Name),
					Field: jsonFieldName,
				}
				changes[field.Name] = changeField
			}
		}
	}
	return changes
}

func getValue(o interface{}, name string) interface{} {
	var r reflect.Value
	if len(name) == 0 {
		return o
	}
	//get the type of the interface provided
	t := reflect.TypeOf(o)
	// if pointer, dereference it and then get the field
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		r = reflect.ValueOf(o).Elem().FieldByName(name)
	} else {
		r = reflect.ValueOf(o).FieldByName(name)
	}
	var val interface{}
	//is the value of that field a pointer?
	if r.Kind() == reflect.Ptr || r.Kind() == reflect.Interface {
		val = reflect.Indirect(r)
	}
	val = r.Interface()
	return val
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
	allUpdates.RegisterInput(wires.EffectsUpdate.GetBufferedOutput(10000))
	PublicSubscribeChange.EnableCopyMode()
	//allUpdates.RegisterInput(wires.NotificationUpdate.GetBufferedOutput(10000))
	subscribeUpdateChannel := allUpdates.GetBufferedOutput(10000)
	go func() {
		for updateObj := range subscribeUpdateChannel {
			update, ok := updateObj.(id.Identifiable)
			if !ok {
				continue
			}
			subscribeablesLock.Acquire("detect change")

			changeDetect, exists := subscribeables[update.GetType()+update.GetId()]
			if !exists {
				continue
			}

			changes := getAllChanged(update, changeDetect.changeDetects)
			if len(changes) != 0 {
				changeDetect.TargetOutput <- &ChangeNotify{
					Type:    changeDetect.Type,
					Id:      changeDetect.Id,
					Changes: changes,
				}
			}
			subscribeablesLock.Release()
		}
	}()
}

func getAllChanged(o interface{}, fields map[string]interface{}) []*ChangeField {
	changedFields := make([]*ChangeField, 0)
	for key, value := range fields {
		switch value.(type) {
		case *ChangeField:
			fieldChange := value.(*ChangeField)
			currentValue := getValue(o, key)
			if !reflect.DeepEqual(currentValue, fieldChange.Value) {
				fieldChange.Value = currentValue
				changedFields = append(changedFields, &ChangeField{fieldChange.Field, fieldChange.Value})
			}
		case map[string]interface{}:
			innerObject := getValue(o, key)
			changedFields = append(changedFields, getAllChanged(innerObject, value.(map[string]interface{}))...)
		}
	}
	return changedFields
}

type SubscribeUpdate struct {
	Type          string
	Id            string
	changeDetects map[string]interface{}
	TargetOutput  chan interface{} `json:"-"`
}

type ChangeField struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type ChangeNotify struct {
	Type    string          `json:"type"`
	Id      string          `json:"uuid"`
	Changes []*ChangeField  `json:"changes"`
	Object  id.Identifiable `json:"-"`
}

func (cn *ChangeNotify) GetId() string {
	return cn.Id
}

func (cn *ChangeNotify) GetType() string {
	return "Change-Notify"
}

//
//func GetCurrentValues() []*ChangeNotify {
//	subscribeablesLock.Acquire("current values")
//	defer subscribeablesLock.Release()
//	values := make([]*ChangeNotify, 0)
//	for _, value := range subscribeables {
//		currentVals := make([]*ChangeField, 0)
//		for _, val := range value.changeDetects {
//			currentVals = append(currentVals, val)
//		}
//
//		newVal := &ChangeNotify{
//			Type:    value.Type,
//			Id:      value.Id,
//			Changes: currentVals,
//		}
//		values = append(values, newVal)
//	}
//	return values
//}

func Test() {

}
