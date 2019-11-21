package objects

import (
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"
)

type Portfolio struct {
	UserUUID string `json:"user_uuid"`
	Uuid     string `json:"uuid"`
	Wallet   int64  `json:"wallet" change:"-"`
	NetWorth int64  `json:"net_worth" change:"-"`
	Level    int64  `json:"level" change:"-"`
}

func (port Portfolio) GetId() string {
	return port.Uuid
}

const PortIdentifiableType = "portfolio"

func (Portfolio) GetType() string {
	return PortIdentifiableType
}

type Effect struct {
	PortfolioUuid string         `json:"portfolio_uuid"`
	Uuid          string         `json:"uuid"`
	Title         string         `json:"title" change:"-"`
	Duration      utils.Duration `json:"duration"`
	StartTime     time.Time      `json:"time"`
	Type          string         `json:"type"`
	InnerEffect   interface{}    `json:"-" change:"inner"`
	Tag           string         `json:"tag"`
}

const EffectIdentifiableType = "effect"

func (Effect) GetType() string {
	return EffectIdentifiableType
}

func (e Effect) GetId() string {
	return e.Uuid
}

type Stock struct {
	Uuid           string        `json:"uuid"`
	Name           string        `json:"name"`
	TickerId       string        `json:"ticker_id"`
	CurrentPrice   int64         `json:"current_price" change:"-"`
	OpenShares     int64         `json:"open_shares" change:"-"`
	ChangeDuration time.Duration `json:"-"`
}

const StockIdentifiableType = "stock"

func (stock Stock) GetId() string {
	return stock.Uuid
}
func (stock Stock) GetType() string {
	return StockIdentifiableType
}

type User struct {
	UserName      string                 `json:"-"`
	Password      string                 `json:"-"`
	DisplayName   string                 `json:"display_name" change:"-"`
	Uuid          string                 `json:"-"`
	PortfolioId   string                 `json:"portfolio_uuid"`
	Active        bool                   `json:"active" change:"-"`
	Config        map[string]interface{} `json:"-"`
	ConfigStr     string                 `json:"-"`
	ActiveClients int64                  `json:"-"`
}

func (user User) GetId() string {
	return user.Uuid
}
func (user User) GetType() string {
	return "user"
}

type Item struct {
	Uuid            string      `json:"uuid"`
	Name            string      `json:"name"`
	ConfigId        string      `json:"config"`
	Type            string      `json:"type"`
	PortfolioUuid   string      `json:"portfolio_uuid"`
	CreateTime      time.Time   `json:"create_time"`
	InnerItemString string      `json:"-"`
	InnerItem       interface{} `json:"-" change:"inner"`
}

func (i Item) GetId() string {
	return i.Uuid
}

const ItemIdentifiableType = "item"

func (Item) GetType() string {
	return ItemIdentifiableType
}

type Ledger struct {
	Uuid         string `json:"uuid"`
	PortfolioId  string `json:"portfolio_id"`
	StockId      string `json:"stock_id"`
	Amount       int64  `json:"amount" change:"-"`
	RecordBookId string `json:"record_book"`
}

const LedgerIdentifiableType = "ledger"

func (ledger Ledger) GetId() string {
	return ledger.Uuid
}

func (Ledger) GetType() string {
	return LedgerIdentifiableType
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

const RecordIdentifiableType = "record_entry"

func (Record) GetType() string {
	return RecordIdentifiableType
}
func (br Record) GetId() string {
	return br.Uuid
}

type Book struct {
	Uuid          string            `json:"uuid"`
	LedgerUuid    string            `json:"ledger_uuid"`
	PortfolioUuid string            `json:"portfolio_uuid"`
	ActiveRecords []ActiveBuyRecord `json:"active_records" change:"-"`
}

type ActiveBuyRecord struct {
	RecordUuid string `json:"record_uud"`
	AmountLeft int64  `json:"still_own"`
}

const BookIdentifiableType = "record_book"

func (Book) GetType() string {
	return BookIdentifiableType
}

func (b Book) GetId() string {
	return b.Uuid
}

type Notification struct {
	Uuid          string      `json:"uuid"`
	PortfolioUuid string      `json:"portfolio_uuid"`
	Timestamp     time.Time   `json:"time"`
	Type          string      `json:"type"`
	Notification  interface{} `json:"notification"`
	Seen          bool        `json:"seen"`
}

func (note Notification) GetId() string {
	return note.Uuid
}

const NotificationIdentifiableType = "notification"

func (Notification) GetType() string {
	return NotificationIdentifiableType
}
