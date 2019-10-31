package record

import (
	"fmt"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/sender"

	"github.com/ThisWillGoWell/stock-simulator-server/src/change"

	"github.com/ThisWillGoWell/stock-simulator-server/src/wires"

	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"

	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
)

var recordsLock = lock.NewLock("records")
var books = make(map[string]*Book)
var records = make(map[string]*Record)
var portfolioBooks = make(map[string][]*Book)

const EntryIdentifiableType = "record_entry"
const BookIdentifiableType = "record_book"
const BuyRecordType = "buy"
const SellRecordType = "sell"

//type Record interface {
//	GetType() string
//	GetId() string
//	GetTime() time.Time
//	GetRecordType() string
//}

type Book struct {
	Uuid             string            `json:"uuid"`
	LedgerUuid       string            `json:"ledger_uuid"`
	PortfolioUuid    string            `json:"portfolio_uuid"`
	ActiveBuyRecords []ActiveBuyRecord `json:"buy_records" change:"-"`
}

type ActiveBuyRecord struct {
	RecordUuid string
	AmountLeft int64
}

type Record struct {
	Uuid           string    `json:"uuid"`
	SharePrice     int64     `json:"share_price"`
	ShareCount     int64     `json:"share_count"`
	Time           time.Time `json:"time"`
	RecordBookUuid string    `json:"book_uuid"`
	Fees           int64     `json:"fee"`
	Taxes          int64     `json:"taxes"`
	Bonus          int64     `json:"bonus"`
	Result         int64     `json:"result"`
}

//func (br *BuyRecord) GetTime() time.Time {
//	return br.Time
//}
//func (*BuyRecord) GetRecordType() string {
//	return BuyRecordType
//}

//type SellRecord struct {
//	Uuid       string `json:"uuid"`
//	SharePrice int64  `json:"share_price"`
//	ShareCount     int64  `json:"amount"`
//}

func NewRecord(recordBookUuid string, amount, sharePrice, taxes, fees, bonus, result int64) {
	uuid := utils.SerialUuid()
	r := MakeRecord(uuid, recordBookUuid, amount, sharePrice, taxes, fees, bonus, result, time.Now())
	wires.RecordsNewObject.Offer(r)
}

func DeleteRecordBook(uuid string) {
	recordsLock.Acquire("delete-record")
	defer recordsLock.Release()
	b, ok := books[uuid]
	if !ok {
		log.Log.Warnf("got delete for record book that we dont know uuid=%s", uuid)
		return
	}
	delete(books, uuid)
	if _, ok := portfolioBooks[b.PortfolioUuid]; ok {
		remove := -1
		for i, l := range portfolioBooks[b.PortfolioUuid] {
			if l.Uuid == uuid {
				remove = 1
				break
			}
		}
		if remove != -1 {
			portfolioBooks[b.PortfolioUuid][remove] = portfolioBooks[b.PortfolioUuid][len(portfolioBooks[b.PortfolioUuid])-1]
		}
		portfolioBooks[b.PortfolioUuid] = portfolioBooks[b.PortfolioUuid][:len(portfolioBooks[b.PortfolioUuid])-1]
	}
}

func MakeBook(uuid, ledgerUuid, portfolioUuid string) error {

	book := &Book{
		Uuid:             uuid,
		LedgerUuid:       ledgerUuid,
		PortfolioUuid:    portfolioUuid,
		ActiveBuyRecords: make([]ActiveBuyRecord, 0),
	}
	bookChange := make(chan interface{})
	if err := change.RegisterPrivateChangeDetect(book, bookChange); err != nil {
		return err
	}
	books[uuid] = book
	if _, ok := portfolioBooks[portfolioUuid]; !ok {
		portfolioBooks[portfolioUuid] = make([]*Book, 0)
	}
	portfolioBooks[portfolioUuid] = append(portfolioBooks[portfolioUuid], books[uuid])

	sender.RegisterChangeUpdate(portfolioUuid, bookChange)
	sender.SendNewObject(portfolioUuid, books[uuid])
	utils.RegisterUuid(uuid, books[uuid])
	return nil
}

func MakeRecord(uuid, recordBookUuid string, amount, sharePrice, taxes, fees, bonus, result int64, t time.Time) (*Record, error) {
	recordsLock.Acquire("new-record")
	defer recordsLock.Release()

	book, ok := books[recordBookUuid]
	if !ok {
		panic("record book not found " + recordBookUuid)
	}
	newRecord := &Record{
		Uuid:           uuid,
		SharePrice:     sharePrice,
		Time:           t,
		ShareCount:     amount,
		RecordBookUuid: recordBookUuid,
		Fees:           fees,
		Bonus:          bonus,
		Result:         result,
		Taxes:          taxes,
	}
	records[uuid] = newRecord
	if amount > 0 {
		book.ActiveBuyRecords = append(book.ActiveBuyRecords, ActiveBuyRecord{RecordUuid: uuid, AmountLeft: amount})
	} else {
		walkRecords(book, amount*-1, true)
	}
	utils.RegisterUuid(uuid, newRecord)

	wires.BookUpdate.Offer(book)
	sender.SendNewObject(book.PortfolioUuid, newRecord)
	return newRecord
}

func walkRecords(book *Book, shares int64, mark bool) int64 {
	amountCleared := 0
	lastAmountCleared := int64(0)
	sharesLeft := shares
	totalCost := int64(0)
	for sharesLeft != 0 {
		if amountCleared >= len(book.ActiveBuyRecords) {
			fmt.Println("WRONG")
		}
		lastAmountCleared = sharesLeft
		activeBuyRecord := book.ActiveBuyRecords[amountCleared]
		record := records[activeBuyRecord.RecordUuid]
		removedShares := activeBuyRecord.AmountLeft

		if activeBuyRecord.AmountLeft > sharesLeft {
			removedShares = sharesLeft
			sharesLeft = 0
		} else if activeBuyRecord.AmountLeft == sharesLeft {
			lastAmountCleared = 0
			amountCleared += 1
			sharesLeft = 0
		} else {
			sharesLeft = sharesLeft - activeBuyRecord.AmountLeft
			amountCleared += 1
		}
		totalCost += removedShares * record.SharePrice

	}

	if mark {
		book.ActiveBuyRecords = book.ActiveBuyRecords[amountCleared:]
		if len(book.ActiveBuyRecords) != 0 {
			book.ActiveBuyRecords[0].AmountLeft -= lastAmountCleared
		}
	}
	return totalCost
}

func GetPrinciple(recordUuid string, shares int64) int64 {
	book := books[recordUuid]
	recordsLock.Acquire("get-principle")
	defer recordsLock.Release()
	return walkRecords(book, shares, false)
}

func GetRecordsForPortfolio(portfolioUuid string) ([]*Book, []*Record) {
	recordsLock.Acquire("get-records")
	defer recordsLock.Release()
	books := portfolioBooks[portfolioUuid]
	portRecord := make([]*Record, 0)

	for _, b := range books {
		for _, active := range b.ActiveBuyRecords {
			portRecord = append(portRecord, records[active.RecordUuid])
		}
	}
	return books, portRecord
}

func GetAllBooks() []*Book {
	recordsLock.Acquire("get-all-books")
	defer recordsLock.Release()
	bookList := make([]*Book, len(books))
	i := 0
	for _, book := range books {
		bookList[i] = book
		i += 1
	}
	return bookList
}
func GetAllRecords() []*Record {
	recordsLock.Acquire("get-all-books")
	defer recordsLock.Release()
	recordList := make([]*Record, len(records))
	i := 0
	for _, record := range records {
		recordList[i] = record
		i += 1
	}
	return recordList
}

func (*Record) GetType() string {
	return EntryIdentifiableType
}
func (br *Record) GetId() string {
	return br.Uuid
}

func (*Book) GetType() string {
	return BookIdentifiableType
}

func (b *Book) GetId() string {
	return b.Uuid
}
