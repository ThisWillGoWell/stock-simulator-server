package database

import (
	"database/sql"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/ledger"
)

var (
	ledgerTableName            = `ledger`
	ledgerTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + ledgerTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`portfolio_id text NOT NULL,` +
		`stock_id text NOT NULL,` +
		`amount bigint NOT NULL,` +
		`record_id text NOT NULL, ` +
		`PRIMARY KEY(uuid)` +
		`);`

	ledgerTableUpdateInsert = `INSERT into ` + ledgerTableName + `(uuid, portfolio_id, record_id, stock_id, amount ) values($1, $2, $3, $4, $5) ` +
		`ON CONFLICT (uuid) DO UPDATE SET amount=EXCLUDED.amount`

	ledgerTableQueryStatement = "SELECT uuid, portfolio_id, stock_id, record_id,  amount FROM " + ledgerTableName + `;`

	ledgerHistoryTableName            = `ledger_history`
	ledgerHistoryTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + ledgerHistoryTableName +
		`( ` +
		`time TIMESTAMPTZ NOT NULL,` +
		`uuid text NOT NULL,` +
		`portfolio_id text NOT NULL, ` +
		`stock_id text NOT NULL, ` +
		`amount bigint NULL` +
		`);`

	ledgerHistoryTableUpdateInsert = `INSERT INTO ` + ledgerHistoryTableName + `(time, uuid, portfolio_id, stock_id, amount) values (NOW(),$1, $2, $3, $4)`

	validLedgerFields = map[string]bool{
		"amount": true,
	}
)

func (d *Database) InitLedger() error {
	if err := d.Exec("ledgers-init", ledgerTableCreateStatement); err != nil {
		return err
	}
	return d.Exec("ledgers-history-init", ledgerHistoryTableCreateStatement)
}

func (d *Database) WriteLedger(entry *ledger.Entry) error {
	if err := d.Exec("ledger-update", ledgerTableUpdateInsert, entry.Uuid, entry.PortfolioId, entry.RecordBookId, entry.StockId, entry.Amount); err != nil {
		return err
	}
	return d.Exec("ledger-history", ledgerHistoryTableUpdateInsert, entry.Uuid, entry.PortfolioId, entry.StockId, entry.Amount)
}

func (d *Database) MakeLedgerHistoryTimeQuery(uuid, timeLength, field, intervalLength string) ([][]interface{}, error) {
	if _, valid := validLedgerFields[field]; !valid {
		return nil, fmt.Errorf("not valid choice")
	}
	return MakeHistoryTimeQuery(ledgerHistoryTableName, uuid, timeLength, field, intervalLength)
}

func (d *Database) MakeLedgerHistoryLimitQuery(uuid, field string, limit int) ([][]interface{}, error) {
	if _, valid := validLedgerFields[field]; !valid {
		return nil, fmt.Errorf("not valid choice")
	}
	return MakeHistoryLimitQuery(ledgerHistoryTableName, uuid, field, limit)
}

func (d *Database) populateLedger() error {
	var uuid, portfolioId, stockId, recordId string
	var amount int64

	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(ledgerTableQueryStatement); err != nil {
		return fmt.Errorf("failed to query portfolio err=%v", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		if err = rows.Scan(&uuid, &portfolioId, &stockId, &recordId, &amount); err != nil {
			return err
		}
		if _, err = ledger.MakeLedgerEntry(uuid, portfolioId, stockId, recordId, amount); err != nil {
			return err
		}
	}
	return rows.Err()
}
