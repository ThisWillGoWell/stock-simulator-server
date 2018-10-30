package metrics

import (
	"time"
)

const interval = time.Second * 10

type MeticCounters struct {
	Object       *ObjectMetric     `json:"object"`
	Update       *ObjectMetric     `json:"update"`
	Connectivity *ConnectiveMetric `json:"connectivity"`
}

var Counter = MeticCounters{
	Object:       ObjectCounter,
	Update:       UpdateCounter,
	Connectivity: ConnectiveCounter,
}

func RunMetrics() {
	runObjectMetrics()

	go func() {
		for {
			select {
			case <-time.After(interval):
				Counter.Update.markInterval()
				Counter.Object.markInterval()
				Counter.Connectivity.markInterval()
			}
		}
	}()

}
