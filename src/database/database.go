package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"

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

func InitDatabase(enableDb, enableDbWrite bool, host, port, username, password string) error {
	db := &Database{}
	if !enableDb {
		db.enable = false
		return nil
	}
	connectionString := fmt.Sprintf()
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

	runHistoricalQueries()
	if !disableDbWrite {
		fmt.Println("starting db writer")
		go databaseWriter()
	}

}

func (d *Database) Execute(writes []interface{}, deletes []interface{}) error {
	if !d.enable {
		return nil
	}
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction err=[%v]", err)
	}
	if writes != nil {
		for _, obj := range writes {
			switch obj.(type) {
			case models.Portfolio:
				err = writePortfolio(obj.(models.Portfolio), tx)
			case models.Ledger:
				err = writeLedger(obj.(models.Ledger), tx)
			case models.User:
				err = writeUser(obj.(models.User), tx)
			case models.Item:
				err = writeItem(obj.(models.Item), tx)
			case models.Stock:
				err = writeStock(obj.(models.Stock), tx)
			case models.Effect:
				err = writeEffect(obj.(models.Effect), tx)
			case models.Record:
				err = writeRecord(obj.(models.Record), tx)
			}
			if err != nil {
				err := fmt.Errorf("failed to write %d items, failed on %T err=[%v]", len(writes), obj, err)
				log.Log.Error(err)
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Log.Errorf("failed to rollback during a failed Exec")
				}
				return err
			}
		}
	}
	if deletes != nil {
		for _, obj := range deletes {
			switch obj.(type) {
			case models.Portfolio:
				err = deletePortfolio(obj.(models.Portfolio), tx)
			case models.Ledger:
				err = deleteLedger(obj.(models.Ledger), tx)
			case models.User:
				err = deleteUser(obj.(models.User), tx)
			case models.Item:
				err = deleteItem(obj.(models.Item), tx)
			case models.Stock:
				err = deleteStock(obj.(models.Stock), tx)
			case models.Effect:
				err = deleteEffect(obj.(models.Effect), tx)
			case models.Record:
				err = deleteRecord(obj.(models.Record), tx)
			case models.Notification:
				err = deleteNotification(obj.(models.Notification), tx)
			}
			if err != nil {
				err := fmt.Errorf("failed to write %d items, failed on %T err=[%v]", len(deletes), obj, err)
				log.Log.Error(err)
				if rollbackErr := tx.Rollback(); rollbackErr != nil {
					log.Log.Errorf("failed to rollback during a failed Exec")
				}
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Log.Errorf("failed to rollback during a failed commit err=[%v]", err)
		}
		return fmt.Errorf("failed to commit data")
	}
	return nil
}

func (d *Database) Exec(commandName, exec string, args ...interface{}) error {
	if !d.enable {
		return nil
	}
	tx, err := d.db.Begin()
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
