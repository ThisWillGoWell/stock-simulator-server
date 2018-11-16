package utils

import (
	"encoding/json"
	"fmt"
)

func PrintJson(e interface{}) {
	val, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(val))
}
