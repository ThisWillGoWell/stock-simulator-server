package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/stock-simulator-server/src/lock"
	"os"
)

var db *sql.DB
var dbLock = lock.NewLock("db lock")

func InitDatabase() {
	conStr := os.Getenv("DB_URI")
	// if the env is not set, default to use the local host default port
	database, err := sql.Open("postgres", conStr)
	if err != nil {
		panic("could not connect to database: " + err.Error())
	}

	db = database
	initLedger()
	initStocks()
	initPortfolio()

	//populateLedger()
	populateStocks()
	//populatePortfolios()

	runLedgerUpdate()
	runStockUpdate()
	runPortfolioUpdate()

}
