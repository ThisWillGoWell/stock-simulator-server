package database

import (
	"log"
	"time"

	"github.com/stock-simulator-server/src/valuable"
)

var (
	stocksTableName            = `stocks`
	stocksTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + stocksTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`ticker_id text NOT NULL,` +
		`name text NOT NULL,` +
		`current_price int,` +
		`open_shares int,` +
		`change_interval numeric(16, 4), ` +
		`PRIMARY KEY(uuid)` +
		`);`

	stocksTableUpdateInsert = `INSERT into ` + stocksTableName + `(uuid, ticker_id, name, current_price, open_shares, change_interval) values($1, $2, $3, $4, $5, $6) ` +
		`ON CONFLICT (uuid) DO UPDATE SET current_price=EXCLUDED.current_price, open_shares=EXCLUDED.open_shares`

	stocksTableQueryStatement = "SELECT uuid, ticker_id, name, current_price, open_shares, change_interval FROM " + stocksTableName
	//getCurrentPrice()z
)

func initStocks() {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		panic("could not begin stocks init: " + err.Error())
	}
	_, err = tx.Exec(stocksTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating metrics table " + err.Error())
	}
	tx.Commit()
}

func writeStock(stock *valuable.Stock) {
	dbLock.Acquire("update-stock")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin stocks init: " + err.Error())
	}
	_, err = tx.Exec(stocksTableUpdateInsert, stock.Uuid, stock.TickerId, stock.Name, stock.CurrentPrice, stock.OpenShares, stock.ChangeDuration)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert stock in table " + err.Error())
	}
	tx.Commit()
}

func populateStocks() {
	var uuid, name, tickerId string
	var currentPrice, openShares int64
	var changeInterval float64

	rows, err := db.Query(stocksTableQueryStatement)
	if err != nil {
		log.Fatal("error quiering databse", err)
		panic("could not populate portfolios: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &tickerId, &name, &currentPrice, &openShares, &changeInterval)
		if err != nil {
			panic(err)
			log.Fatal(err)
		}
		t := time.Duration(changeInterval)
		valuable.MakeStock(uuid, tickerId, name, currentPrice, openShares, t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
