package ledger

import (
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/lock"
)

var Entries = make(map[string]*Entry)
var EntriesLock = lock.NewLock("ledger-entries-lock")

type Entry struct {
	Uuid string `uuid:"uuid"`
	ExchangeId string `json:"exchange_id"`
	PortfolioId string `json:"portfolio_id"`
	Amount        float64 `json:"amount"`
	updateChannel duplicator.ChannelDuplicator `json:"-"`
}

func NewLedgerEntry( portfolioId, exchangeId string, amount float64 ){
	EntriesLock.Acquire("new ledger entry")
	defer EntriesLock.Release()
	uuid := utils.PseudoUuid()
	for{
		if _ ,exists := Entries[uuid]; !exists{
			break
		}
		uuid = utils.PseudoUuid()
	}
	Entries[uuid] = &Entry{
		entries
	}
}