package ledger

import (
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"

	"github.com/ThisWillGoWell/stock-simulator-server/src/database"

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
	models.Ledger
	Lock          *lock.Lock                    `json:"-"`
	UpdateChannel *duplicator.ChannelDuplicator `json:"-"`
}

/**
build a new ledger entry and generate a new uuid for it
takes in the lock acquired since trade already owns the lock for the entries
*/
func NewLedgerEntry(portfolioId, stockId string) (*Entry, error) {
	EntriesLock.Acquire("new-entry")
	defer EntriesLock.Release()

	uuid := utils.SerialUuid()
	recordId := utils.SerialUuid()

	e, err := MakeLedgerEntry(uuid, portfolioId, stockId, recordId, 0, true)
	if err != nil {
		return nil, err
	}

	if err = database.Db.WriteLedger(e); err != nil {
		_ = DeleteLedger(uuid, true)
		return nil, err
	}

	wires.LedgerNewObject.Offer(e)
	return e, nil
}

func DeleteLedger(uuid string, lockAquired bool) error {
	if !lockAquired {
		EntriesLock.Acquire("delete-ledger")
		defer EntriesLock.Release()
	}
	var l *Entry
	var ok bool
	if l, ok = Entries[uuid]; !ok {
		log.Log.Warnf("go delete for uuid not found uuid=%s", uuid)
		return nil
	}
	var err error
	err = database.Db.DeleteLedger(uuid)

	delete(Entries, uuid)
	if _, ok := EntriesPortfolioStock[l.PortfolioId]; ok {
		delete(EntriesPortfolioStock[l.PortfolioId], uuid)
	}
	if _, ok := EntriesStockPortfolio[l.PortfolioId]; ok {
		delete(EntriesStockPortfolio[l.PortfolioId], uuid)
	}

	change.UnregisterChangeDetect(l)
	l.UpdateChannel.StopDuplicator()
	utils.RemoveUuid(uuid)
	return err
}

func (l *Entry) PublishUpdate() error {
	if err := database.Db.WriteLedger(l.Ledger); err != nil {
		return err
	}
	l.UpdateChannel.Offer(l)
	return nil
}

/**
Make a Ledger
*/
func MakeLedgerEntry(uuid, portfolioId, stockId, recordId string, amount int64, lockAquired bool) (*Entry, error) {
	if !lockAquired {
		EntriesLock.Acquire("make-ledger")
		defer EntriesLock.Release()
	}
	entry := &Entry{
		Ledger: models.Ledger{
			Uuid:         uuid,
			PortfolioId:  portfolioId,
			StockId:      stockId,
			Amount:       amount,
			RecordBookId: recordId,
		},
		UpdateChannel: duplicator.MakeDuplicator(fmt.Sprintf("LedgerEntry-%s", uuid)),
	}
	if err := record.MakeBook(recordId, uuid, portfolioId); err != nil {
		return nil, fmt.Errorf("failed to make recored book for ledger err=[%v]", err)
	}

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
