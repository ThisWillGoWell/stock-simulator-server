package record

import (
	"time"

	"github.com/stock-simulator-server/src/change"

	"github.com/stock-simulator-server/src/wires"

	"github.com/stock-simulator-server/src/lock"

	"github.com/stock-simulator-server/src/utils"
)

var recordsLock = lock.NewLock("records")
var books = make(map[string]*Book)
var records = make(map[string]*Record)

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
	MakeRecord(uuid, recordBookUuid, amount, sharePrice, taxes, fees, bonus, result, time.Now())
	//sender.SendNewObject(portfolioUuid, record)
}

func MakeBook(uuid, ledgerUuid string) {

	books[uuid] = &Book{
		Uuid:             uuid,
		LedgerUuid:       ledgerUuid,
		ActiveBuyRecords: make([]ActiveBuyRecord, 0),
	}
	change.RegisterPublicChangeDetect(books[uuid])
}

func MakeRecord(uuid, recordBookUuid string, amount, sharePrice, taxes, fees, bonus, result int64, t time.Time) *Record {
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
	wires.RecordsNewObject.Offer(newRecord)
	wires.BookUpdate.Offer(book)
	return newRecord
}

func walkRecords(book *Book, shares int64, mark bool) int64 {
	amountCleared := 0
	lastAmountCleared := int64(0)
	sharesLeft := shares
	totalCost := int64(0)

	for sharesLeft != 0 {
		lastAmountCleared = sharesLeft
		activeBuyRecord := book.ActiveBuyRecords[amountCleared]
		record := records[activeBuyRecord.RecordUuid]
		removedShares := activeBuyRecord.AmountLeft

		if activeBuyRecord.AmountLeft >= sharesLeft {
			removedShares = sharesLeft
			sharesLeft = 0
		} else {
			sharesLeft = sharesLeft - activeBuyRecord.AmountLeft
		}
		totalCost += removedShares * record.SharePrice

		if sharesLeft != 0 {
			amountCleared += 1
		}
	}
	if sharesLeft == 0 {
		amountCleared += 1
	}
	if mark {
		book.ActiveBuyRecords[0].AmountLeft -= lastAmountCleared
		book.ActiveBuyRecords = book.ActiveBuyRecords[amountCleared:]
	}
	return totalCost
}

func GetPrinciple(recordUuid string, shares int64) int64 {
	book := books[recordUuid]
	recordsLock.Acquire("get-principle")
	defer recordsLock.Release()
	return walkRecords(book, shares, false)
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
