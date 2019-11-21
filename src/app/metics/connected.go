package metrics

import "sync/atomic"

type ConnectiveInnerMetrics struct {
	NumberClients int64  `json:"num_clients"`
	TxCount       uint64 `json:"tx"`
	RxCount       uint64 `json:"rx"`
	TxSize        uint64 `json:"tx_size"`
	RxSize        uint64 `json:"rx_size"`
}

type ConnectiveMetric struct {
	Total        ConnectiveInnerMetrics `json:"total"`
	Interval     ConnectiveInnerMetrics `json:"interval"`
	lastInterval ConnectiveInnerMetrics `json:"-"`
}

var ConnectiveCounter = &ConnectiveMetric{}

func (cc *ConnectiveMetric) markInterval() {
	cc.Interval.TxCount = cc.Total.TxCount - cc.lastInterval.TxCount
	cc.Interval.RxCount = cc.Total.RxCount - cc.lastInterval.RxCount
	cc.Interval.NumberClients = cc.Total.NumberClients - cc.lastInterval.NumberClients
	cc.Interval.RxSize = cc.Total.RxSize - cc.lastInterval.RxSize
	cc.Interval.TxSize = cc.Total.TxSize - cc.lastInterval.TxSize
	cc.lastInterval = cc.Total
}

func SendMessage(size int) {
	atomic.AddUint64(&ConnectiveCounter.Total.TxSize, uint64(size))
	atomic.AddUint64(&ConnectiveCounter.Total.TxCount, 1)
}
func RecieveMessage(size int) {
	atomic.AddUint64(&ConnectiveCounter.Total.RxSize, uint64(size))
	atomic.AddUint64(&ConnectiveCounter.Total.RxCount, 1)

}

func ClientConnect() {
	atomic.AddInt64(&ConnectiveCounter.Total.NumberClients, 1)

}
func ClientDisconnect() {
	atomic.AddInt64(&ConnectiveCounter.Total.NumberClients, -1)

}
