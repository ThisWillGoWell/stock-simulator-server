package database

import (
	"github.com/stock-simulator-server/src/valuable"
	"log"
	"time"
)

var (
	stocksTableName            = `stocks`
	stocksTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + stocksTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`ticker_id text NOT NULL,` +
		`name text NOT NULL,` +
		`current_price numeric(16, 4),` +
		`open_shares int,` +
		`PRIMARY KEY(uuid)` +
		`);`

	stocksTableUpdateInsert = `INSERT into ` + stocksTableName + `(uuid, ticker_id, name, current_price, open_shares) values($1, $2, $3, $4, $5) ` +
		`ON CONFLICT (uuid) DO UPDATE SET current_price=EXCLUDED.current_price, open_shares=EXCLUDED.open_shares`

	stocksTableQueryStatement = "SELECT * FROM " + stocksTableName
	//getCurrentPrice()
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

func runStockUpdate() {
	stockUpdateChannel := valuable.ValuableUpdateChannel.GetBufferedOutput(100)
	go func() {
		for stockUpdated := range stockUpdateChannel {
			stock := stockUpdated.(*valuable.Stock)
			updateStock(stock)
		}
	}()
}

func updateStock(stock *valuable.Stock) {
	dbLock.Acquire("update-stock")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin stocks init")
	}
	_, err = tx.Exec(stocksTableUpdateInsert, stock.Uuid, stock.TickerId, stock.Name, stock.CurrentPrice, stock.OpenShares)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert stock in table " + err.Error())
	}
	tx.Commit()
}

func populateStocks() {
	var uuid, name, tickerId string
	var currentPrice, openShares float64
	var databaseId int

	rows, err := db.Query(stocksTableQueryStatement)
	if err != nil {
		log.Fatal("error quiering databse")
		panic("could not populate portfolios: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&databaseId, &uuid, &tickerId, &name, &currentPrice, &openShares)
		if err != nil {
			panic(err)
			log.Fatal(err)
		}
		valuable.MakeStock(uuid, tickerId, name, currentPrice, openShares, time.Second*60)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
