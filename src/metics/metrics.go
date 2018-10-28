package metics

import (
	"encoding/json"
	"fmt"
	"time"
)

const interval = time.Second * 10

type MeticCounters struct {
	Object       *ObjectMetric      `json:"object"`
	Update       *ObjectMetric      `json:"update"`
	Duplicator   *DuplicatorMetrics `json"duplocator"`
	Connectivity *ConnectiveMetrics `json"connective"`
}

var Counter = MeticCounters{
	Object:     ObjectCounter,
	Update:     UpdateCounter,
	Duplicator: DuplicatorCounter,
}

func RunMetrics() {
	runObjectMetrics()

	go func() {
		for {
			select {
			case <-time.After(interval):
				Counter.Update.markInterval()
				Counter.Object.markInterval()
				Counter.Duplicator.markInterval()
				b, _ := json.Marshal(Counter)
				str := string(b)
				fmt.Println(str)

			}
		}
	}()

}
