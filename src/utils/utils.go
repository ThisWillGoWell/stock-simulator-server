package utils

import (
	"reflect"
	"strings"
	"fmt"
	"encoding/json"
	"time"
)

/*
	So this is some of the most bs code I have written in awhile
	So how can we detect the change in a pointer's value in any struct?
	welp this works
	The idea is to use tags of a struct to watch changes on a value using reflect
	using the tag "change" we can tag a value to be watched for changes

	type changeTest struct {
			Id	string
			Test int `change:"Id"`
	}

	Here the change detect is registered and will be tagged with the Id Field and its value

	&changeTest{
		Id: "this is my id",
		Test: 1,
	}

	[{
		"tags": {
			"Id": "this is my id"
		},
		"name": "Test",
		"value": 2
	}]




 */


const(
	changeTag     = "change"
	)

var(
	//subscribeables is something that can be subscribed to
	subscribeables = make(map[interface{}]*SubscribeUpdate)
	subscribeablesLock = NewLock("subscribeables")
	SubscribeUpdates = MakeDuplicator()
)

type Identifiable interface {
	Id() string
	GetType() string
}

func registerChangeDetect(o Identifiable)(*SubscribeUpdate){
	//get the include tags
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface{
		t = reflect.ValueOf(o).Elem().Type()
	}
	//t := reflect.ValueOf(o).Type()
	newChangeDetect := &SubscribeUpdate{
		changeDetects: make(map[string]*ChangeField),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value, exists := field.Tag.Lookup(changeTag)

		if exists{
			tagList := strings.Split(value, ",")
			value := getValue(o, field.Name)
			fmt.Println("postGet", value, reflect.TypeOf(value))
			tags := make(map[string]interface{})
			for _, tag := range tagList{
				 tags[tag] = getValue(o, tag)
			}
			changeField := &ChangeField{
				Tag: tags,
				Value: value,
				Name: field.Name,
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

//todo this is not idea, sould try and refrence these right off the pointer?
func updateIds(o interface{}, changeField *ChangeField){
	for tag :=range changeField.Tag{
		changeField.Tag[tag] = getValue(o, tag)
	}

}

func StartDetectChanges(){
	subscribeUpdateChannel := SubscribeUpdates.GetOutput()
	go func(){
		for update := range subscribeUpdateChannel {
			subscribeablesLock.Acquire("detect change")

			changeDetect, exists := subscribeables[update]
			if ! exists{
				changeDetect = registerChangeDetect(update)
				subscribeables[update] = changeDetect
				fmt.Println(changeDetect.changeDetects)
			}
			changedFields := make([]*ChangeField, 0)
			changed := false
			for filedName, fieldChange := range changeDetect.changeDetects{
				currentValue := getValue(update, filedName)
				fmt.Println(currentValue, fieldChange.Value, currentValue == fieldChange.Value)
				if currentValue != fieldChange.Value {
					changed = true
					fmt.Println("change!")
					fieldChange.Value = currentValue
					updateIds(update, fieldChange)
					changedFields = append(changedFields, fieldChange)
				}
			}
			if changed{
				str, _ := json.Marshal(changedFields)
				fmt.Println(string(str))
			}
			subscribeablesLock.Release()
		}
	}()

}

type SubscribeUpdate struct {
	changeDetects map[string]*ChangeField
	}

type ChangeField struct {
	Tag map[string]interface{} `json:"tag"`
	Name string `json:"name"`
	Value interface{} `json:"value"`
	}

type changeTest struct {
	Id	string
	Test int `change:"Id"`
}

func Test(){
	StartDetectChanges()
	SubscribeUpdates.EnableCopyMode()

	v := &changeTest{
		Id: "this is my id",
		Test: 1,
	}


	tempUpdateChannel := make(chan interface{})
	SubscribeUpdates.RegisterInput(tempUpdateChannel)
	SubscribeUpdates.Offer(v)
	SubscribeUpdates.Offer(v)
	SubscribeUpdates.Offer(v)
	v.Test = 2
	SubscribeUpdates.Offer(v)
	v.Test = 3
	SubscribeUpdates.Offer(v)
	for{
		time.Sleep(10 * time.Hour)
	}
}


type Subscribe interface {

}
