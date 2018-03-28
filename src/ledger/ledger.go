package ledger

import (
	"fmt"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/utils"
)

const objectType = "ledger"

//map of uuid -> entry
var Entries = make(map[string]*Entry)
var EntriesStockPortfolio = make(map[string]map[string]*Entry)
var EntriesPortfolioStock = make(map[string]map[string]*Entry)

// map of stock_uuid -> open shares
var EntriesLock = lock.NewLock("ledger-entries-lock")
var EntriesUpdate = duplicator.MakeDuplicator("ledger-entries-update")

type Entry struct {
	Lock          *lock.Lock                    `json:"-"`
	Uuid          string                        `uuid:"uuid"`
	PortfolioId   string                        `json:"portfolio_id"`
	StockId       string                        `json:"portfolio_id"`
	Amount        float64                       `json:"amount" change:"-"`
	UpdateChannel *duplicator.ChannelDuplicator `json:"-"`
}

func NewLedgerEntry(portfolioId, stockId string, amount float64) *Entry {
	EntriesLock.Acquire("make ledger entry")
	defer EntriesLock.Release()
	uuid := utils.PseudoUuid()
	for {
		if _, exists := Entries[uuid]; !exists {
			break
		}
		uuid = utils.PseudoUuid()
	}
	return MakeLedgerEntry(uuid, portfolioId, stockId, 0)
}

func MakeLedgerEntry(uuid, portfolioId, stockId string, amount float64) *Entry {

	entry := &Entry{
		Uuid:          uuid,
		PortfolioId:   portfolioId,
		Amount:        amount,
		StockId:       stockId,
		UpdateChannel: duplicator.MakeDuplicator(fmt.Sprintf("LedgerEntry-%s", uuid)),
	}

	Entries[uuid] = entry
	EntriesPortfolioStock[portfolioId][stockId] = entry
	EntriesStockPortfolio[stockId][portfolioId] = entry
	EntriesUpdate.RegisterInput(entry.UpdateChannel.GetOutput())
	EntriesUpdate.Offer(entry)
	return entry
}

func RemoveLedgerEntry(uuid string) {
	entry := Entries[uuid]

	delete(EntriesStockPortfolio[entry.StockId], entry.PortfolioId)
	if len(EntriesStockPortfolio[entry.StockId]) == 0 {
		delete(EntriesStockPortfolio, entry.StockId)
	}

	delete(EntriesPortfolioStock[entry.PortfolioId], entry.StockId)
	if len(EntriesPortfolioStock[entry.PortfolioId]) == 0 {
		delete(EntriesPortfolioStock, entry.PortfolioId)
	}
	delete(Entries, entry.Uuid)
	duplicator.UnlinkDouplicator(EntriesUpdate, entry.UpdateChannel)
}

func GetAllLedgers() []*Entry {
	EntriesLock.Acquire("get-all-ledgers")
	defer EntriesLock.Release()
	lst := make([]*Entry, len(Entries))
	i := 0
	for _, val := range Entries {
		lst[i] = val
		i += 1
	}
	return lst
}

func (ledger *Entry) GetId() string {
	return ledger.Uuid
}

func (ledger *Entry) GetType() string {
	return objectType
}
