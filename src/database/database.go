package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stock-simulator-server/src/lock"
	"os"
	"time"
)

var db *sql.DB
var ts *sql.DB

var dbLock = lock.NewLock("db lock")

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
	//populateLedger()
	//populateStocks()
	//populatePortfolios()

	runLedgerUpdate()
	runStockUpdate()
	runPortfolioUpdate()
	runStockHistoryUpdate()
	runPortfolioHistoryUpdate()
}
