package database

import (
	"github.com/stock-simulator-server/src/valuable"
)

var (
	stocksHistoryTableName            = `stocks_history`
	stocksHistoryTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + stocksHistoryTableName +
		`( ` +
		`time TIMESTAMPTZ NOT NULL,` +
		`uuid text NOT NULL,` +
		`current_price numeric(16, 4) NULL,` +
		`open_shares int NULL` +
		`);`
	stocksHistoryTSInit = `CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE; SELECT create_hypertable('` + stocksHistoryTableName + `', 'time');`

	stocksHistoryTableUpdateInsert = `INSERT INTO ` + stocksHistoryTableName + `(time, uuid, current_price, open_shares) values (NOW(),$1, $2, $3);`

	stocksHistroyTableQueryStatement = "SELECT * FROM " + stocksHistoryTableName + " WHERE uuid=$1"
	//getCurrentPrice()
)

func initStocksHistory() {
	tx, err := ts.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin stocks init: " + err.Error())
	}
	_, err = tx.Exec(stocksHistoryTableCreateStatement)
	if err != nil {

	}
	tx.Commit()
	tx, err = ts.Begin()
	_, err = tx.Exec(stocksHistoryTSInit)
	if err != nil {

	}
	tx.Commit()
}

func runStockHistoryUpdate() {
	newStockChannel := valuable.NewStockChannel.GetBufferedOutput(10)
	stockUpdateChannel := valuable.ValuableUpdateChannel.GetBufferedOutput(100)

	go func() {
		for stockNew := range newStockChannel {
			stock := stockNew.(*valuable.Stock)
			updateStockHistory(stock)
		}
	}()

	go func() {
		for stockUpdated := range stockUpdateChannel {
			stock := stockUpdated.(*valuable.Stock)
			updateStockHistory(stock)
		}
	}()
}

func updateStockHistory(stock *valuable.Stock) {
	tx, err := ts.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin stocks init")
	}
	_, err = tx.Exec(stocksHistoryTableUpdateInsert, stock.Uuid, stock.CurrentPrice, stock.OpenShares)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert stock in table " + err.Error())
	}
	tx.Commit()
}
