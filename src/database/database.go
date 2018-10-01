package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/change"
	"github.com/stock-simulator-server/src/duplicator"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/lock"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/valuable"
)

var db *sql.DB
var ts *sql.DB

var dbLock = lock.NewLock("db lock")

var DatabseWriter = duplicator.MakeDuplicator("database-writer")

func InitDatabase() {
	dbConStr := os.Getenv("DB_URI")
	// if the env is not set, default to use the local host default port
	database, err := sql.Open("postgres", dbConStr)
	fmt.Println(dbConStr)
	if err != nil {
		panic("could not connect to database: " + err.Error())
	}
	db = database

	for i := 0; i < 10; i++ {
		err := db.Ping()

		if err == nil {
			break
		}
		fmt.Println("waitng for connection to db")
		<-time.After(time.Second)
	}

	conStr := os.Getenv("TS_URI")
	timeseriers, err := sql.Open("postgres", conStr)
	if err != nil {
		panic("could not connect to database: " + err.Error())
	}

	ts = timeseriers

	for i := 0; i < 10; i++ {
		err := timeseriers.Ping()
		if err == nil {
			break
		}
		fmt.Println("waitng for connection to ts")
		<-time.After(time.Second)
	}

	initLedger()
	initStocks()
	initPortfolio()
	initStocksHistory()
	initPortfolioHistory()

	initLedgerHistory()
	initAccount()

	populateLedger()
	populateStocks()
	populatePortfolios()
	populateUsers()

	for _, l := range ledger.Entries {
		port := portfolio.Portfolios[l.PortfolioId]
		stock := valuable.Stocks[l.StockId]
		port.UpdateInput.RegisterInput(stock.UpdateChannel.GetBufferedOutput(10))
		port.UpdateInput.RegisterInput(l.UpdateChannel.GetBufferedOutput(10))

	}
	for _, port := range portfolio.Portfolios {
		port.Update()
	}

	runHistoricalQueries()
	go databaseWriter()

}

func databaseWriter() {
	write := DatabseWriter.GetBufferedOutput(1000)
	for obj := range write {
		// pull from uuid map since change-detect come across as change,
		//todo fix race condition from update and pulling current val
		val, exists := utils.GetVal(obj.(change.Identifiable).GetId())
		if !exists {
			panic("db write for uuid not in uuidmap: " + obj.(change.Identifiable).GetId())
		}
		switch val.(type) {
		case *portfolio.Portfolio:
			writePortfolio(val.(*portfolio.Portfolio))
			writePortfolioHistory(val.(*portfolio.Portfolio))
		case *account.User:
			writeUser(val.(*account.User))
		case *ledger.Entry:
			writeLedger(val.(*ledger.Entry))
			writeLedgerHistory(val.(*ledger.Entry))
		case *valuable.Stock:
			writeStock(val.(*valuable.Stock))
			writeStockHistory(val.(*valuable.Stock))
		default:
			panic("deafult call of databse writer")
		}
	}
}
