package utils

import (
	"crypto/rand"
	"fmt"
)

func PseudoUuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
		return
	}

	uuid = fmt.Sprintf("%X", b[0:7])

	return uuid
}
