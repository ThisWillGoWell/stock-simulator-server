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
		`wallet numeric(16, 4) NOT NULL,` +
		`PRIMARY KEY(uuid)` +
		`);`

	portfolioTableUpdateInsert = `INSERT into ` + portfolioTableName + `(uuid, name, wallet) values($1, $2, $3) ` +
		`ON CONFLICT (uuid) DO UPDATE SET wallet=EXCLUDED.wallet;`

	portfolioTableQueryStatement = "SELECT * FROM " + portfolioTableName + `;`
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
		panic("error occurred while creating metrics table " + err.Error())
	}
	tx.Commit()
}

func runPortfolioUpdate() {
	portfolioUpdateChannel := portfolio.PortfoliosUpdateChannel.GetBufferedOutput(100)
	go func() {
		for portfolioUpdated := range portfolioUpdateChannel {
			port := portfolioUpdated.(*portfolio.Portfolio)
			updatePortfolio(port)
		}
	}()

}

func updatePortfolio(port *portfolio.Portfolio) {
	dbLock.Acquire("update-stock")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin stocks init")
	}
	_, err = tx.Exec(portfolioTableUpdateInsert, port.UUID, port.Name, port.Wallet)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert stock in table " + err.Error())
	}
	tx.Commit()
}

func populatePortfolios() {
	var uuid, name string
	var wallet float64

	rows, err := db.Query(portfolioTableQueryStatement)
	if err != nil {
		log.Fatal("error quiering databse")
		panic("could not populate portfolios: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &name, &wallet)
		if err != nil {
			log.Fatal(err)
		}
		portfolio.MakePortfolio(uuid, name, wallet)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
