package ledger

import (
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/record"

	"github.com/ThisWillGoWell/stock-simulator-server/src/change"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"github.com/ThisWillGoWell/stock-simulator-server/src/duplicator"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
)

const objectType = "ledger"

//map of uuid -> entry
var Entries = make(map[string]*Entry)
var EntriesStockPortfolio = make(map[string]map[string]*Entry)
var EntriesPortfolioStock = make(map[string]map[string]*Entry)

// map of stock_uuid -> open shares
var EntriesLock = lock.NewLock("ledger-entries-lock")

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
	Amount        int64                         `json:"amount" change:"-"`
	UpdateChannel *duplicator.ChannelDuplicator `json:"-"`
	RecordBookId  string                        `json:"record_book"`
}

/**
build a new ledger entry and generate a new uuid for it
takes in the lock acquired since trade already owns the lock for the entries
*/
func NewLedgerEntry(portfolioId, stockId string, lockAcquired bool) (*Entry, error) {
	if !lockAcquired {
		EntriesLock.Acquire("make ledger entry")
		defer EntriesLock.Release()
	}
	uuid := utils.SerialUuid()
	recordId := utils.SerialUuid()

	return MakeLedgerEntry(uuid, portfolioId, stockId, recordId, 0)

}

func deleteLedger(uuid string) {
	// does not delete the item in the database as this is only called
	// when a ledger fails to create
	EntriesLock.Acquire("delete-ledger")
	defer EntriesLock.Release()
	var l *Entry
	var ok bool
	if l, ok = Entries[uuid]; !ok {
		log.Log.Warnf("go delete for uuid not found uuid=%s", uuid)
		return
	}
	delete(Entries, uuid)
	if _, ok := EntriesPortfolioStock[l.PortfolioId]; ok {
		delete(EntriesPortfolioStock[l.PortfolioId], uuid)
	}
	if _, ok := EntriesStockPortfolio[l.PortfolioId]; ok {
		delete(EntriesStockPortfolio[l.PortfolioId], uuid)
	}
	change.UnregisterChangeDetect(l)
	utils.RemoveUuid(uuid)
}

/**
Make a Ledger
*/
func MakeLedgerEntry(uuid, portfolioId, stockId, recordId string, amount int64) (*Entry, error) {
	fmt.Println("making ledger ")
	entry := &Entry{
		Uuid:          uuid,
		PortfolioId:   portfolioId,
		Amount:        amount,
		StockId:       stockId,
		RecordBookId:  recordId,
		UpdateChannel: duplicator.MakeDuplicator(fmt.Sprintf("LedgerEntry-%s", uuid)),
	}
	record.MakeBook(recordId, uuid, portfolioId)

	if err := change.RegisterPublicChangeDetect(entry); err != nil {
		return nil, err
	}

	if EntriesPortfolioStock[portfolioId] == nil {
		EntriesPortfolioStock[portfolioId] = make(map[string]*Entry)
	}
	EntriesPortfolioStock[portfolioId][stockId] = entry

	if EntriesStockPortfolio[stockId] == nil {
		EntriesStockPortfolio[stockId] = make(map[string]*Entry)
	}
	Entries[uuid] = entry
	EntriesStockPortfolio[stockId][portfolioId] = entry
	entry.UpdateChannel.EnableCopyMode()

	wires.LedgerNewObject.Offer(entry)
	wires.LedgerUpdate.RegisterInput(entry.UpdateChannel.GetOutput())
	utils.RegisterUuid(uuid, entry)
	return entry, nil
}

/*
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
