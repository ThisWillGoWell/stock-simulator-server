package valuable

import (
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/wires/duplicator"
)

var ValuablesLock = lock.NewLock("valuables")

/**
Valuable was an old thing that i used to abstract stocks one more level but just became cumbersome
*/
type Valuable interface {
	GetId() string
	GetName() string
	GetValue() float64
	GetLock() *lock.Lock
	GetUpdateChannel() *duplicator.ChannelDuplicator
	Update()
}
