package database

import (
	"database/sql"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"

	"github.com/ThisWillGoWell/stock-simulator-server/src/valuable"
	"github.com/pkg/errors"
)

var (
	stocksTableName            = `stocks`
	stocksTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + stocksTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`ticker_id text NOT NULL,` +
		`name text NOT NULL,` +
		`current_price bigint,` +
		`open_shares bigint,` +
		`change_interval numeric(16, 4), ` +
		`PRIMARY KEY(uuid)` +
		`);`

	stocksTableUpdateInsert = `INSERT into ` + stocksTableName + `(uuid, ticker_id, name, current_price, open_shares, change_interval) values($1, $2, $3, $4, $5, $6) ` +
		`ON CONFLICT (uuid) DO UPDATE SET current_price=EXCLUDED.current_price, open_shares=EXCLUDED.open_shares`

	stocksTableQueryStatement = "SELECT uuid, ticker_id, name, current_price, open_shares, change_interval FROM " + stocksTableName

	stocksHistoryTableName            = `stocks_history`
	stocksHistoryTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + stocksHistoryTableName +
		`( ` +
		`time TIMESTAMPTZ NOT NULL,` +
		`uuid text NOT NULL,` +
		`current_price bigint NULL,` +
		`open_shares bigint NULL` +
		`);`

	stocksHistoryTableUpdateInsert = `INSERT INTO ` + stocksHistoryTableName + `(time, uuid, current_price, open_shares) values (NOW(),$1, $2, $3);`

	stocksTableDeleteStatement        = "DELETE from " + stocksTableName + `WHERE uuid = $1`
	stocksHistoryTableDeleteStatement = "DELETE from " + stocksHistoryTableName + `WHERE uuid = $1`

	validStockFields = map[string]bool{
		"current_price": true,
	}
)

func (d *Database) InitStocks() error {
	if err := d.Exec("stocks-init", stocksTableCreateStatement); err != nil {
		return err
	}
	return d.Exec("stocks-history-init", stocksHistoryTableCreateStatement)
}

func (d *Database) WriteStock(stock models.Stock) error {
	if err := d.Exec(stocksTableUpdateInsert, stock.Uuid, stock.TickerId, stock.Name, stock.CurrentPrice, stock.OpenShares, stock.ChangeDuration); err != nil {
		return err
	}
	return d.Exec(stocksHistoryTableUpdateInsert, stock.Uuid, stock.CurrentPrice, stock.OpenShares)

}

func (d *Database) DeleteStock(uuid string) error {
	e1 := d.Exec("delete-stock", stocksTableDeleteStatement, uuid)
	e2 := d.Exec("delete-stock", stocksHistoryTableDeleteStatement, uuid)
	if e1 != nil || e2 != nil {
		return fmt.Errorf("delete stock uuid=%s stockdb=[%v] history=[%v]", uuid, e1, e2)
	}
	return nil
}

func (d *Database) MakeStockHistoryTimeQuery(uuid, timeLength, field, intervalLength string) ([][]interface{}, error) {
	if _, valid := validStockFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return d.MakeHistoryTimeQuery(stocksHistoryTableName, uuid, timeLength, field, intervalLength)

}

func (d *Database) MakeStockHistoryLimitQuery(uuid, field string, limit int) ([][]interface{}, error) {
	if _, valid := validStockFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return d.MakeHistoryLimitQuery(stocksHistoryTableName, uuid, field, limit)
}

func (d *Database) populateStocks() error {
	var uuid, name, tickerId string
	var currentPrice, openShares int64
	var changeInterval float64

	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(stocksTableQueryStatement); err != nil {
		return fmt.Errorf("failed to query portfolio err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		if err = rows.Scan(&uuid, &tickerId, &name, &currentPrice, &openShares, &changeInterval); err != nil {
			return err
		}
		if _, err = valuable.MakeStock(uuid, tickerId, name, currentPrice, openShares, t); err != nil {
			return err
		}
	}
	return rows.Err()
}
