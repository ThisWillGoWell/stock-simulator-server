package utils

import (
	"fmt"
	"time"

	"github.com/stock-simulator-server/src/lock"
)

var reclaimedUuid = make([]string, 0)
var uuidMap = make(map[string]interface{})
var uuidLock = lock.NewLock("uuid-lock")
var counterNum = 0

var deleteDelta = time.Minute * 2

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
	if len(reclaimedUuid) != 0 {
		uuid = reclaimedUuid[0]
		reclaimedUuid = reclaimedUuid[1:]
	} else {
		for {
			counterNum += 1
			if _, exists := uuidMap[uuid]; !exists {
				break
			}
			uuid = fmt.Sprintf("%d", counterNum)
		}
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
	go func() {
		fmt.Println("got delete uuid for: " + uuid)
		<-time.After(deleteDelta)
		fmt.Println("ready to delete uuid: " + uuid)
		uuidLock.Acquire("reclaim-uuid")
		defer uuidLock.Release()
		delete(uuidMap, uuid)
		reclaimedUuid = append(reclaimedUuid, uuid)
		fmt.Println("uuid deleted")
	}()
}
