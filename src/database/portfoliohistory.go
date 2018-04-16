package database

import (
	"github.com/stock-simulator-server/src/portfolio"
)

var (
	portfolioHistoryTableName            = `portfolio_history`
	portfolioHistoryTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + portfolioHistoryTableName +
		`( ` +
		`time TIMESTAMPTZ NOT NULL,` +
		`uuid text NOT NULL,` +
		`net_worth numeric(16, 4) NULL,` +
		`wallet numeric(16, 4) NULL` +
		`);`
	portfolioHistoryTSInit = `CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE; SELECT create_hypertable('` + portfolioHistoryTableName + `', 'time');`

	portfolioHistoryTableUpdateInsert = `INSERT INTO ` + portfolioHistoryTableName + `(time, uuid, net_worth, wallet) values (NOW(),$1, $2, $3)`

	portfolioHistroyTableQueryStatement = "SELECT * FROM " + portfolioHistoryTableName + " WHERE uuid"
	//getCurrentPrice()
)

func initPortfolioHistory() {
	tx, err := ts.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin portfolio init: " + err.Error())
	}
	_, err = tx.Exec(portfolioHistoryTableCreateStatement)
	if err != nil {

	}
	tx.Commit()
	tx, err = ts.Begin()
	_, err = tx.Exec(portfolioHistoryTSInit)
	if err != nil {

	}
	tx.Commit()
}

func runPortfolioHistoryUpdate() {
	portfolioUpdateChannel := portfolio.PortfoliosUpdateChannel.GetBufferedOutput(100)
	portfolioNewChannel := portfolio.NewPortfolioChannel.GetBufferedOutput(10)

	go func() {
		for portfolioNew := range portfolioNewChannel {
			port := portfolioNew.(*portfolio.Portfolio)
			updatePortfolioHistory(port)
		}
	}()

	go func() {
		for portfolioUpdate := range portfolioUpdateChannel {
			port := portfolioUpdate.(*portfolio.Portfolio)
			updatePortfolioHistory(port)
		}
	}()
}

func updatePortfolioHistory(port *portfolio.Portfolio) {
	tx, err := ts.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin portfolio init")
	}
	_, err = tx.Exec(portfolioHistoryTableUpdateInsert, port.UUID, port.NetWorth, port.Wallet)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert portfolio in table " + err.Error())
	}
	tx.Commit()
}
