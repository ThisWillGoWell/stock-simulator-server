package valuable

import (
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/duplicator"
)

var Valuables = make(map[string]Valuable)
var ValuablesLock = lock.NewLock("valuables")
var ValuableUpdateChannel = duplicator.MakeDuplicator()

type Valuable interface {
	GetId() string
	GetValue() float64
	GetLock() *lock.Lock
	GetUpdateChannel() *duplicator.ChannelDuplicator
}
