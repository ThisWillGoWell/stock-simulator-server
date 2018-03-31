package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/stock-simulator-server/src/lock"
)

var uuidMap = make(map[string]bool)
var uuidLock = lock.NewLock("uuid-lock")

func PseudoUuid() string {
	uuidLock.Acquire("new uuid")
	defer uuidLock.Release()
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	uuid := fmt.Sprintf("%X", b[0:7])
	for {
		if _, exists := uuidMap[uuid]; !exists {
			break
		}
		_, err = rand.Read(b)
		if err != nil {
			panic(err)
		}
		uuid = fmt.Sprintf("%X", b[0:7])
	}
	uuidMap[uuid] = true
	return uuid
}

func RegisterUuid(uuid string) {
	uuidLock.Acquire("new uuid")
	defer uuidLock.Release()
	uuidMap[uuid] = true
}
