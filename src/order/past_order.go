package order

import "time"

type OrderRecord struct {
	Order *Order
	time.Time
	LedgerUuid string
}
