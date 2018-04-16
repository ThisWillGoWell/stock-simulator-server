package ledger

import (
	"fmt"
	"github.com/stock-simulator-server/src/change"
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

/**
Ledgers store who owns what stock
They are all done though uuid strings since that's all that's required
They are stored in two maps
1) given a stock uuid, get all portfolios that own it
2) given a portfolio uuid, get all stocks it owns
*/
type Entry struct {
	Lock          *lock.Lock                    `json:"-"`
	Uuid          string                        `json:"uuid"`
	PortfolioId   string                        `json:"portfolio_id"`
	StockId       string                        `json:"stock_id"`
	Amount        float64                       `json:"amount" change:"-"`
	UpdateChannel *duplicator.ChannelDuplicator `json:"-"`
}

/**
build a new ledger entry and generate a new uuid for it
takes in the lock acquired since trade already owns the lock for the entries
*/
func NewLedgerEntry(portfolioId, stockId string, lockAcquired bool) *Entry {
	if !lockAcquired {
		EntriesLock.Acquire("make ledger entry")
		defer EntriesLock.Release()
	}
	uuid := utils.PseudoUuid()

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
	if EntriesPortfolioStock[portfolioId] == nil {
		EntriesPortfolioStock[portfolioId] = make(map[string]*Entry)
	}
	EntriesPortfolioStock[portfolioId][stockId] = entry

	if EntriesStockPortfolio[stockId] == nil {
		EntriesStockPortfolio[stockId] = make(map[string]*Entry)
	}
	EntriesStockPortfolio[stockId][portfolioId] = entry
	entry.UpdateChannel.EnableCopyMode()

	change.NewSubscribeCreated.Offer(entry)
	EntriesUpdate.RegisterInput(entry.UpdateChannel.GetOutput())
	utils.RegisterUuid(uuid, entry)
	return entry
}

/**
remove a ledger form both maps
*/
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

/**
get All ledgers so they can be sent on connection
*/
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
