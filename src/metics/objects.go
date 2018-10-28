package metics

import (
	"sync/atomic"

	"github.com/stock-simulator-server/src/wires"
)

type InnerObjectMetric struct {
	Portfolio     uint64 `json:"portfolio"`
	Stocks        uint64 `json:"stocks"`
	Ledger        uint64 `json:"ledger"`
	Users         uint64 `json:"users"`
	Records       uint64 `json:"records"`
	Notifications uint64 `json:"notifications"`
	Items         uint64 `json:"items"`
}

type ObjectMetric struct {
	Total        InnerObjectMetric `json:"total"`
	Interval     InnerObjectMetric `json:"interval"`
	LastInterval InnerObjectMetric `json:"last_interval"`
}

var ObjectCounter = &ObjectMetric{}
var UpdateCounter = &ObjectMetric{}

func (om *ObjectMetric) markInterval() {
	om.Interval.Items = om.Total.Items - om.LastInterval.Items
	om.Interval.Portfolio = om.Total.Portfolio - om.LastInterval.Portfolio
	om.Interval.Stocks = om.Total.Stocks - om.LastInterval.Stocks
	om.Interval.Ledger = om.Total.Ledger - om.LastInterval.Ledger
	om.Interval.Users = om.Total.Users - om.LastInterval.Users
	om.Interval.Records = om.Total.Records - om.LastInterval.Records
	om.Interval.Notifications = om.Total.Notifications - om.LastInterval.Notifications
	om.LastInterval = om.Total
}

func runObjectMetrics() {
	go func() {
		for range wires.PortfolioNewObject.GetBufferedOutput(10000) {
			atomic.AddUint64(&ObjectCounter.Total.Portfolio, 1)
		}
	}()
	go func() {
		for range wires.PortfolioUpdate.GetBufferedOutput(10000) {
			atomic.AddUint64(&UpdateCounter.Total.Portfolio, 1)
		}
	}()

	go func() {
		for range wires.LedgerNewObject.GetBufferedOutput(10000) {
			atomic.AddUint64(&ObjectCounter.Total.Ledger, 1)
		}
	}()
	go func() {
		for range wires.LedgerUpdate.GetBufferedOutput(10000) {
			atomic.AddUint64(&UpdateCounter.Total.Ledger, 1)
		}
	}()

	go func() {
		for range wires.UsersNewObject.GetBufferedOutput(10000) {
			atomic.AddUint64(&ObjectCounter.Total.Users, 1)
		}
	}()
	go func() {
		for range wires.UsersUpdate.GetBufferedOutput(10000) {
			atomic.AddUint64(&UpdateCounter.Total.Users, 1)
		}
	}()
	//note: records do not have an update
	go func() {
		for range wires.RecordsNewObject.GetBufferedOutput(10000) {
			atomic.AddUint64(&ObjectCounter.Total.Records, 1)
		}
	}()

	go func() {
		for range wires.NotificationNewObject.GetBufferedOutput(10000) {
			atomic.AddUint64(&ObjectCounter.Total.Notifications, 1)
		}
	}()
	go func() {
		for range wires.NotificationUpdate.GetBufferedOutput(10000) {
			atomic.AddUint64(&UpdateCounter.Total.Notifications, 1)
		}
	}()

	go func() {
		for range wires.ItemsNewObjects.GetBufferedOutput(10000) {
			atomic.AddUint64(&ObjectCounter.Total.Items, 1)
		}
	}()
	go func() {
		for range wires.ItemsUpdate.GetBufferedOutput(10000) {
			atomic.AddUint64(&UpdateCounter.Total.Items, 1)
		}
	}()
	go func() {
		for range wires.StocksNewObject.GetBufferedOutput(10000) {
			atomic.AddUint64(&ObjectCounter.Total.Stocks, 1)
		}
	}()
	go func() {
		for range wires.StocksUpdate.GetBufferedOutput(10000) {
			atomic.AddUint64(&UpdateCounter.Total.Stocks, 1)
		}
	}()

}
