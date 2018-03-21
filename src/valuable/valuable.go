package valuable

import (
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/duplicator"
)

var ValuablesLock = lock.NewLock("valuables")
var ValuableUpdateChannel = duplicator.MakeDuplicator("valueable-update")
var NewValuableChannel = duplicator.MakeDuplicator("new-vauleuable")

type Valuable interface {
	GetId() string
	GetName() string
	GetValue() float64
	GetLock() *lock.Lock
	GetUpdateChannel() *duplicator.ChannelDuplicator
}
