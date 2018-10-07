package database

import (
	"log"

	"github.com/stock-simulator-server/src/portfolio"
)

var (
	portfolioTableName            = `portfolio`
	portfolioTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + portfolioTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`name text NOT NULL,` +
		`bigint int NOT NULL,` +
		`level int NOT NULL, ` +
		`PRIMARY KEY(uuid) ` +
		`);`

	portfolioTableUpdateInsert = `INSERT into ` + portfolioTableName + `(uuid, name, wallet, level) values($1, $2, $3, $4) ` +
		`ON CONFLICT (uuid) DO UPDATE SET wallet=EXCLUDED.wallet,  level=EXCLUDED.level;`

	portfolioTableQueryStatement = "SELECT uuid, name, wallet, level FROM " + portfolioTableName + `;`
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
	dbLock.Acquire("update-portfolio")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin portfolio init: " + err.Error())
	}
	_, err = tx.Exec(portfolioTableUpdateInsert, port.Uuid, port.UserUUID, port.Wallet, port.Level)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert portfolio in table " + err.Error())
	}
	tx.Commit()
}

func populatePortfolios() {
	var uuid, userUuid string
	var wallet, level int64

	rows, err := db.Query(portfolioTableQueryStatement)
	if err != nil {
		log.Fatal("error query database", err)
		panic("could not populate portfolios: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &userUuid, &wallet, &level)
		if err != nil {
			log.Fatal(err)
		}
		portfolio.MakePortfolio(uuid, userUuid, wallet, level)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
