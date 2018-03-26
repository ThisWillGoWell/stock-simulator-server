package database

import (
	"github.com/stock-simulator-server/src/portfolio"
	"log"
)


var(
	ledger = `portfolio`
	portfolioTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + portfolioTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`name text NOT NULL,`+
		`wallet numeric(16, 4) NOT NULL,` +
		`PRIMARY KEY(uuid)` +
		`);`

	portfolioTableUpdateInsert = `INSERT into ` + portfolioTableName + `(uuid, name, wallet, net_worth) values($1, $2, $3, $4) `+
		`ON CONFLICT (uuid) DO UPDATE SET wallet=EXCLUDED.wallet, net_worth=EXCLUDED.net_worth`

	portfolioTableQueryStatement = "SELECT * FROM " + portfolioTableName + `;`
	//getCurrentPrice()
)

func initPortfolio(){
	tx, err := db.Begin()
	if err != nil{
		db.Close()
		panic("could not begin stocks init: " + err.Error())
	}
	_,err = tx.Exec(portfolioTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating metrics table " + err.Error())
	}
	tx.Commit()
}

func runPortfolioUpdate(){
	portfolioUpdateChannel := portfolio.PortfoliosUpdateChannel.GetBufferedOutput(100)
	go func(){
		for portfolioUpdated := range portfolioUpdateChannel{
			port := portfolioUpdated.(*portfolio.Portfolio)
			updatePortfolio(port)
		}
	}();


}

func updatePortfolio(port *portfolio.Portfolio) {
	dbLock.Acquire("update-stock")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin stocks init")
	}
	_, err = tx.Exec(portfolioTableUpdateInsert, port.UUID, port.Name, port.Wallet, port.NetWorth)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert stock in table " + err.Error())
	}
	tx.Commit()
}

func populatePortfolios(){
	var uuid, name string
	var wallet float64

	rows, err := db.Query(portfolioTableQueryStatement)
	if err != nil{
		log.Fatal("error quiering databse")
		panic("could not populate portfolios: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		uuid, name string

		err := rows.Scan(&loadedPortfolio.UUID, &loadedPortfolio.Name, &loadedPortfolio.NetWorth)
		if err != nil {
			log.Fatal(err)
		}
		portfolio.MakePortfolio()
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
