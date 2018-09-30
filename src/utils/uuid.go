package utils

import (
	"fmt"
	"github.com/stock-simulator-server/src/lock"
)

var uuidMap = make(map[string]interface{})
var uuidLock = lock.NewLock("uuid-lock")
var counterNum = 0

/**
this is one of the more controversial desgin choices I made
here is the single souce of truth for assigns uuids to objects
it keeps a map of all uuids and a pointer to that uuid so given a uuid, its type can be found
The problem here is uuids don't play well with large scale databases, and it would be better just to
have the database be the one assigning uuid
*/
func SerialUuid() string {
	uuidLock.Acquire("new uuid")
	defer uuidLock.Release()

	uuid := fmt.Sprintf("%d", counterNum)

	for {
		counterNum += 1
		if _, exists := uuidMap[uuid]; !exists {
			break
		}
		uuid = fmt.Sprintf("%d", counterNum)
	}
	uuidMap[uuid] = nil
	return uuid
}

func GetVal(uuid string) (interface{}, bool) {
	uuidLock.Acquire("uuid get")
	defer uuidLock.Release()
	val, exists := uuidMap[uuid]
	return val, exists
}

func RegisterUuid(uuid string, val interface{}) {
	uuidLock.Acquire("uuid register")
	defer uuidLock.Release()
	uuidMap[uuid] = val
}

func RemoveUuid(uuid string) {
	uuidLock.Acquire("free-uuid")
	defer uuidLock.Release()
	if _, exists := uuidMap[uuid]; exists {
		delete(uuidMap, uuid)
	}
}
