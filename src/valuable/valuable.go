package valuable

import "github.com/stock-simulator-server/src/utils"

var Valuables = make(map[string]Valuable)
var ValuablesLock = utils.NewLock("valuables")
var ValuableUpdateChannel = utils.MakeDuplicator()

type Valuable interface {
	GetId() string
	GetValue() float64
	GetLock() *utils.Lock
	GetUpdateChannel() *utils.ChannelDuplicator
}
