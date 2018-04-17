package database

import (
	"github.com/stock-simulator-server/src/ledger"
	"log"
)

var (
	ledgerTableName            = `ledger`
	ledgerTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + ledgerTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`portfolio_id text NOT NULL,` +
		`stock_id text NOT NULL,` +
		`amount numeric(16, 4) NOT NULL,` +
		`PRIMARY KEY(uuid)` +
		`);`

	ledgerTableUpdateInsert = `INSERT into ` + ledgerTableName + `(uuid, portfolio_id, stock_id, amount) values($1, $2, $3, $4) ` +
		`ON CONFLICT (uuid) DO UPDATE SET amount=EXCLUDED.amount`

	ledgerTableQueryStatement = "SELECT * FROM " + ledgerTableName + `;`
	//getCurrentPrice()
)

func initLedger() {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		panic("could not begin ledger init: " + err.Error())
	}
	_, err = tx.Exec(ledgerTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating leger table " + err.Error())
	}
	tx.Commit()
}

func runLedgerUpdate() {
	ledgerUpdateChannel := ledger.EntriesUpdate.GetBufferedOutput(10)
	go func() {
		for portfolioUpdated := range ledgerUpdateChannel {
			entry := portfolioUpdated.(*ledger.Entry)
			updateLedger(entry)
		}
	}()

}

func updateLedger(entry *ledger.Entry) {
	dbLock.Acquire("update-ledger")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin stocks init")
	}
	_, err = tx.Exec(ledgerTableUpdateInsert, entry.Uuid, entry.PortfolioId, entry.StockId, entry.Amount)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert stock in table " + err.Error())
	}
	tx.Commit()
}

func populateLedger() {
	var uuid, portfolioId, stockId string
	var amount float64

	rows, err := db.Query(portfolioTableQueryStatement)
	if err != nil {
		log.Fatal("error quiering databse")
		panic("could not populate portfolios: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &portfolioId, &stockId, amount)
		if err != nil {
			log.Fatal(err)
		}
		ledger.MakeLedgerEntry(uuid, portfolioId, stockId, amount)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
