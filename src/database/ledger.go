package database

import (
	"github.com/ThisWillGoWell/stock-simulator-server/src/ledger"
	"github.com/ThisWillGoWell/stock-simulator-server/src/log"
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
	//getCurrentPrice()
)

func initLedger() {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		log.Alerts.Fatal("could not begin ledger init: ", err.Error())
		panic("could not begin ledger init: " + err.Error())
	}
	_, err = tx.Exec(ledgerTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating leger table " + err.Error())
	}
	tx.Commit()
}

func writeLedger(entry *ledger.Entry) {
	dbLock.Acquire("update-ledger")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin ledger init" + err.Error())
	}
	_, err = tx.Exec(ledgerTableUpdateInsert, entry.Uuid, entry.PortfolioId, entry.RecordBookId, entry.StockId, entry.Amount)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert ledger in table " + err.Error())
	}
	tx.Commit()
}

func populateLedger() {
	var uuid, portfolioId, stockId, recordId string
	var amount int64

	rows, err := db.Query(ledgerTableQueryStatement)
	if err != nil {
		//log.Fatal("error query database", err)
		panic("could not populate portfolios: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &portfolioId, &stockId, &recordId, &amount)
		if err != nil {
			//log.Fatal("error in querying ledger: ", err)
		}
		ledger.MakeLedgerEntry(uuid, portfolioId, stockId, recordId, amount)
	}
	err = rows.Err()
	if err != nil {
		//log.Fatal(err)
	}
}
