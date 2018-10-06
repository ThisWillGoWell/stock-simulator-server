package database

import (
	"log"

	"github.com/stock-simulator-server/src/ledger"
)

var (
	ledgerTableName            = `ledger`
	ledgerTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + ledgerTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`portfolio_id text NOT NULL,` +
		`stock_id text NOT NULL,` +
		`amount int NOT NULL,` +
		`investment_value int NOT NULL, ` +
		`PRIMARY KEY(uuid)` +
		`);`

	ledgerTableUpdateInsert = `INSERT into ` + ledgerTableName + `(uuid, portfolio_id, stock_id, amount, investment_value) values($1, $2, $3, $4, $5) ` +
		`ON CONFLICT (uuid) DO UPDATE SET amount=EXCLUDED.amount, investment_value=EXCLUDED.investment_value`

	ledgerTableQueryStatement = "SELECT uuid, portfolio_id, stock_id, amount, investment_value FROM " + ledgerTableName + `;`
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

func writeLedger(entry *ledger.Entry) {
	dbLock.Acquire("update-ledger")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin ledger init" + err.Error())
	}
	_, err = tx.Exec(ledgerTableUpdateInsert, entry.Uuid, entry.PortfolioId, entry.StockId, entry.Amount, entry.InvestmentValue)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert ledger in table " + err.Error())
	}
	tx.Commit()
}

func populateLedger() {
	var uuid, portfolioId, stockId string
	var amount, investmentVal int64

	rows, err := db.Query(ledgerTableQueryStatement)
	if err != nil {
		log.Fatal("error quiering databse", err)
		panic("could not populate portfolios: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &portfolioId, &stockId, &amount, &investmentVal)
		if err != nil {
			log.Fatal("error in querying ledger: ", err)
		}
		ledger.MakeLedgerEntry(uuid, portfolioId, stockId, amount, investmentVal)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
