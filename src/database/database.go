package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/ledger"
	"github.com/ThisWillGoWell/stock-simulator-server/src/lock"
	"github.com/ThisWillGoWell/stock-simulator-server/src/portfolio"
	"github.com/ThisWillGoWell/stock-simulator-server/src/valuable"
	_ "github.com/lib/pq"
)

var Db *Database

var dbLock = lock.NewLock("db lock")

type Database struct {
	enable bool
	db     sql.DB
}

func InitDatabase(enableDb, enableDbWrite bool, host, port, username, password string) (Database, error) {
	if !enableDb {
		return
	}

	dbConStr := os.Getenv("RDS")
	// if the env is not set, default to use the local host default port
	database, err := sql.Open("postgres", dbConStr)
	fmt.Println(dbConStr)
	if err != nil {
		log.Alerts.Fatal("could not connect to database: " + err.Error())
		log.Log.Fatal("could not connect to database: " + err.Error())
		panic("could not connect to database: " + err.Error())
	}
	db = database

	for i := 0; i < 10; i++ {
		err := db.Ping()

		if err == nil {
			break
		}
		fmt.Println("	waitng for connection to db")
		<-time.After(time.Second)
	}
	log.Log.Println("connected to database")

	initLedger()
	initStocks()
	initPortfolio()
	initStocksHistory()
	initPortfolioHistory()
	initNotification()
	initItems()
	initLedgerHistory()
	initAccount()
	initRecordHistory()
	initEffect()

	populateUsers()
	populateStocks()
	populatePortfolios()
	populateLedger()
	populateItems()
	populateNotification()
	populateRecords()
	populateEffects()

	for _, l := range ledger.Entries {
		port := portfolio.Portfolios[l.PortfolioId]
		stock := valuable.Stocks[l.StockId]
		port.UpdateInput.RegisterInput(stock.UpdateChannel.GetBufferedOutput(100))
		port.UpdateInput.RegisterInput(l.UpdateChannel.GetBufferedOutput(100))
	}
	for _, port := range portfolio.Portfolios {
		port.Update()
	}

	runHistoricalQueries()
	if !disableDbWrite {
		fmt.Println("starting db writer")
		go databaseWriter()
	}

}

func (d *Database) Exec(commandName, exec string, args ...interface{}) error {
	if !d.enable {
		return nil
	}
	tx, err := db.Begin()
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("begin %s: err=[%v]", commandName, err)
	}
	_, err = tx.Exec(itemsTableCreateStatement, args)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("exec %s: command=%v err=[%v]", commandName, exec, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit %s: command=%v err=[%v]", commandName, exec, err)
	}
	return nil
}
