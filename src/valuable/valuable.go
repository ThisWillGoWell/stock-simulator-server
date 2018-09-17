package valuable

import (
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/lock"
)

var ValuablesLock = lock.NewLock("valuables")
var UpdateChannel = duplicator.MakeDuplicator("valuable-update")
var NewObjectChannel = duplicator.MakeDuplicator("new-valuable")

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
