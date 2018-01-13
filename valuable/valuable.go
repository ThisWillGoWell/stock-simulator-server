package valuable

import "github.com/stock-simulator-server/utils"

var Valuables = make(map[string]Valuable)
var ValuablesLock = utils.NewLock("valuables")

type Valuable interface {
	GetID() string
	GetValue() float64
	GetLock() *utils.Lock
	GetUpdateChannel() *utils.ChannelDuplicator
}
