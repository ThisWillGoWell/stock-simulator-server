package database

import (
	"github.com/pkg/errors"
	"github.com/stock-simulator-server/src/valuable"
)

var (
	stocksHistoryTableName            = `stocks_history`
	stocksHistoryTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + stocksHistoryTableName +
		`( ` +
		`time TIMESTAMPTZ NOT NULL,` +
		`uuid text NOT NULL,` +
		`current_price bigint NULL,` +
		`open_shares bigint NULL` +
		`);`
	stocksHistoryTSInit = `CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE; SELECT create_hypertable('` + stocksHistoryTableName + `', 'time');`

	stocksHistoryTableUpdateInsert = `INSERT INTO ` + stocksHistoryTableName + `(time, uuid, current_price, open_shares) values (NOW(),$1, $2, $3);`

	validStockFields = map[string]bool{
		"current_price": true,
	}
	//getCurrentPrice()
)

func initStocksHistory() {
	tx, err := ts.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin stocks history init: " + err.Error())
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

func writeStockHistory(stock *valuable.Stock) {
	tx, err := ts.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin stocks history init: " + err.Error())
	}
	_, err = tx.Exec(stocksHistoryTableUpdateInsert, stock.Uuid, stock.CurrentPrice, stock.OpenShares)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert stock in table " + err.Error())
	}
	tx.Commit()
}

func MakeStockHistoryTimeQuery(uuid, timeLength, field, intervalLength string) ([][]interface{}, error) {
	if _, valid := validStockFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return MakeHistoryTimeQuery(stocksHistoryTableName, uuid, timeLength, field, intervalLength)

}

func MakeStockHistoryLimitQuery(uuid, field string, limit int) ([][]interface{}, error) {
	if _, valid := validStockFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return MakeHistoryLimitQuery(stocksHistoryTableName, uuid, field, limit)
}
