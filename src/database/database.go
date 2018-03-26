package database

import (
	"database/sql"
	"os"
	"github.com/stock-simulator-server/src/lock"
	_ "github.com/lib/pq"

)

var db *sql.DB
var dbLock = lock.NewLock("db lock")

func InitDatabase()  {
	conStr := os.Getenv("DB_URI")
	// if the env is not set, default to use the local host default port
	database, err := sql.Open("postgres", conStr)
	if err != nil{
		panic("could not connect to database: " + err.Error())
	}

	db = database
	initStocks()
	initPortfolio()

	runStockUpdate()
	runPortfolioUpdate()
}


