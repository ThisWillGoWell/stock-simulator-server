package utils

import (
	"reflect"
	"fmt"
	"time"
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


const(
	changeTag     = "change"
	)

var(
	//subscribeables is something that can be subscribed to
	subscribeables        = make(map[string]*SubscribeUpdate)
	subscribeablesLock    = NewLock("subscribeables")
	SubscribeUpdateInputs = MakeDuplicator()
	SubscribeUpdateOutput = MakeDuplicator()
)


func registerChangeDetect(o Identifiable)(*SubscribeUpdate){
	//get the include tags
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface{
		t = reflect.ValueOf(o).Elem().Type()
	}
	//t := reflect.ValueOf(o).Type()
	newChangeDetect := &SubscribeUpdate{
		Type: o.GetType(),
		Id: o.GetId(),
		changeDetects: make(map[string]*ChangeField),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		_, exists := field.Tag.Lookup(changeTag)
		if exists{
			jsonFieldName, exists := field.Tag.Lookup("json")
			if ! exists {
				jsonFieldName = field.Name
			}

			changeField := &ChangeField{
				Value: nil,
				Field: jsonFieldName,
			}
			newChangeDetect.changeDetects[field.Name] = changeField
		}
	}
	return newChangeDetect
}

func getValue(o interface{}, name string) interface{}{
	var r reflect.Value
	//get the type of the interface provided
	t := reflect.TypeOf(o)
	// if pointer, dereference it and then get the field
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface{
		r = reflect.ValueOf(o).Elem().FieldByName(name)
	}else{
		r = reflect.ValueOf(o).FieldByName(name)
	}

	//is the value of that field a pointer?
	if r.Kind() == reflect.Ptr || r.Kind() == reflect.Interface {
		fmt.Println("value")
		return reflect.Indirect(r)
	}
	fmt.Println(o, r.Type(), r.Interface())
	return r.Interface()
}


func StartDetectChanges(){
	subscribeUpdateChannel := SubscribeUpdateInputs.GetOutput()
	go func(){
		for updateObj := range subscribeUpdateChannel {
			update, ok := updateObj.(Identifiable)
			if ! ok {
				 panic("got a non identifiable in the change detector")
			}
			subscribeablesLock.Acquire("detect change")

			changeDetect, exists := subscribeables[update.GetType() + update.GetId()]
			if ! exists{
				changeDetect = registerChangeDetect(update)
				subscribeables[update.GetType() + update.GetId()] = changeDetect
				fmt.Println(changeDetect.changeDetects)
			}
			changedFields := make([]*ChangeField, 0)
			changed := false
			for filedName, fieldChange := range changeDetect.changeDetects{
				currentValue := getValue(update, filedName)

				if ! reflect.DeepEqual(currentValue, fieldChange.Value) {
					changed = true
					fieldChange.Value = currentValue
					changedFields = append(changedFields, fieldChange)
				}
			}
			if changed{
				fmt.Println("changed:", changedFields)
				SubscribeUpdateOutput.Offer(&ChangeNotfity{
					Type: changeDetect.Type,
					Id: changeDetect.Id,
					Changes: changedFields,
				})
			}
			subscribeablesLock.Release()
		}
	}()

}

type SubscribeUpdate struct {
	Type string
	Id string
	changeDetects map[string]*ChangeField
	}


type ChangeField struct {

	Field string `json:"field"`
	Value interface{} `json:"value"`
	}
type ChangeNotfity struct {
	Type string `json:"type"`
	Id string `json:"id"`
	Changes []*ChangeField `json:"changes"`
}

type changeTest struct {
	Name	string
	Test int `change:"-"`
}

func (change *changeTest)GetId() string{
	return change.Name
}

func (change *changeTest)GetType() string{
	return "test"
}

func Test(){
	StartDetectChanges()
	SubscribeUpdateInputs.EnableCopyMode()

	v := &changeTest{
		Name: "this is my name",
		Test: 1,
	}


	tempUpdateChannel := make(chan interface{})
	SubscribeUpdateInputs.RegisterInput(tempUpdateChannel)
	SubscribeUpdateInputs.Offer(v)
	SubscribeUpdateInputs.Offer(v)
	SubscribeUpdateInputs.Offer(v)
	v.Test = 2
	SubscribeUpdateInputs.Offer(v)
	v.Test = 3
	SubscribeUpdateInputs.Offer(v)
	for{
		time.Sleep(10 * time.Hour)
	}
}

