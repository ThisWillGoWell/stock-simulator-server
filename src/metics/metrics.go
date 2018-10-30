package metrics

import (
	"time"
)

const interval = time.Minute * 15

type MetricCounters struct {
	Object       *ObjectMetric     `json:"object"`
	Update       *ObjectMetric     `json:"update"`
	Connectivity *ConnectiveMetric `json:"connectivity"`
}

var Counter = MetricCounters{
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
