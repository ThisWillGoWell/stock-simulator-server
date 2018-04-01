package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stock-simulator-server/src/lock"
	"os"
)

var db *sql.DB
var ts *sql.DB

var dbLock = lock.NewLock("db lock")

func InitDatabase() {
	dbConStr := os.Getenv("DB_URI")
	// if the env is not set, default to use the local host default port
	database, err := sql.Open("postgres", dbConStr)
	if err != nil {
		panic("could not connect to database: " + err.Error())
	}
	db = database

	conStr := os.Getenv("TS_URI")
	timeseriers, err := sql.Open("postgres", conStr)
	if err != nil {
		panic("could not connect to database: " + err.Error())
	}

	ts = timeseriers

	initLedger()
	initStocks()
	initPortfolio()
	initStocksHistory()

	//populateLedger()
	//populateStocks()
	//populatePortfolios()

	runLedgerUpdate()
	runStockUpdate()
	runPortfolioUpdate()
	runStockHistoryUpdate()
}
