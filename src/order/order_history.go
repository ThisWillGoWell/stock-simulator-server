package order

import "time"

type TradeOrderHistory struct {
	StockUuid string
	Time      time.Time
	Portfolio string
	NextOrder Order
}
