package valuable

import "stock-server/utils"

var Valuables = make(map[string]Valuable)
var ValuablesLock = utils.NewLock()

type Valuable interface {
	GetValue() float64
	GetLock() *utils.Lock
	GetUpdateChannel() *utils.ChannelDuplicator
}
