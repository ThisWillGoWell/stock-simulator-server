package metics

import "sync/atomic"

type DuplicatorInnerMetrics struct {
	NumberRunning int64  `json:"num_active"`
	TransferCount uint64 `json:"transfer_count"`
}

type DuplicatorMetrics struct {
	Total        DuplicatorInnerMetrics `json:"total"`
	Interval     DuplicatorInnerMetrics `json:"interval"`
	lastInterval DuplicatorInnerMetrics
}

func (dm *DuplicatorMetrics) markInterval() {
	dm.Interval.NumberRunning = dm.Total.NumberRunning - dm.lastInterval.NumberRunning
	dm.Interval.TransferCount = dm.Total.TransferCount - dm.lastInterval.TransferCount
}

var DuplicatorCounter = &DuplicatorMetrics{}

func SendTransfer() {
	atomic.AddUint64(&DuplicatorCounter.Total.TransferCount, 1)
}

func StartTransfer() {
	atomic.AddInt64(&DuplicatorCounter.Total.NumberRunning, 1)
}
func StopTransfer() {
	atomic.AddInt64(&DuplicatorCounter.Total.NumberRunning, -1)
}
