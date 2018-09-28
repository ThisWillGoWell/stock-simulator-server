package database

import (
	"github.com/stock-simulator-server/src/portfolio"
	"log"
)

var (
	portfolioTableName            = `portfolio`
	portfolioTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + portfolioTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`name text NOT NULL,` +
		`wallet int NOT NULL,` +
		`PRIMARY KEY(uuid)` +
		`);`

	portfolioTableUpdateInsert = `INSERT into ` + portfolioTableName + `(uuid, name, wallet) values($1, $2, $3) ` +
		`ON CONFLICT (uuid) DO UPDATE SET wallet=EXCLUDED.wallet;`

	portfolioTableQueryStatement = "SELECT uuid, name, wallet FROM " + portfolioTableName + `;`
	//getCurrentPrice()
)

func initPortfolio() {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		panic("could not begin stocks init: " + err.Error())
	}
	_, err = tx.Exec(portfolioTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating portfolio table " + err.Error())
	}
	tx.Commit()
}

func writePortfolio(port *portfolio.Portfolio) {
	dbLock.Acquire("update-stock")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin portfolio init: " + err.Error())
	}
	_, err = tx.Exec(portfolioTableUpdateInsert, port.UUID, port.UserUUID, port.Wallet)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert portfolio in table " + err.Error())
	}
	tx.Commit()
}

func populatePortfolios() {
	var uuid, userUuid string
	var wallet int64

	rows, err := db.Query(portfolioTableQueryStatement)
	if err != nil {
		log.Fatal("error quiering databse")
		panic("could not populate portfolios: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &userUuid, &wallet)
		if err != nil {
			log.Fatal(err)
		}
		portfolio.MakePortfolio(uuid, userUuid, wallet)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
