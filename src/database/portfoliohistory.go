package database

import (
	"errors"

	"github.com/stock-simulator-server/src/portfolio"
)

var (
	portfolioHistoryTableName            = `portfolio_history`
	portfolioHistoryTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + portfolioHistoryTableName +
		`( ` +
		`time TIMESTAMPTZ NOT NULL,` +
		`uuid text NOT NULL,` +
		`net_worth bigint NULL,` +
		`wallet bigint NULL` +
		`);`
	portfolioHistoryTSInit = `CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE; SELECT create_hypertable('` + portfolioHistoryTableName + `', 'time');`

	portfolioHistoryTableUpdateInsert = `INSERT INTO ` + portfolioHistoryTableName + `(time, uuid, net_worth, wallet) values (NOW(),$1, $2, $3)`

	//getCurrentPrice()
	validPortfolioFields = map[string]bool{
		"wallet":    true,
		"net_worth": true,
	}
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

func writePortfolioHistory(port *portfolio.Portfolio) {
	tx, err := ts.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin portfolio history write init: " + err.Error())
	}
	_, err = tx.Exec(portfolioHistoryTableUpdateInsert, port.Uuid, port.NetWorth, port.Wallet)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert portfolio in table " + err.Error())
	}
	tx.Commit()
}
func MakePortfolioHistoryTimeQuery(uuid, timeLength, field, intervalLength string) ([][]interface{}, error) {
	if _, valid := validPortfolioFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return MakeHistoryTimeQuery(portfolioHistoryTableName, uuid, timeLength, field, intervalLength)

}

func MakePortfolioHistoryLimitQuery(uuid, field string, limit int) ([][]interface{}, error) {
	if _, valid := validPortfolioFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return MakeHistoryLimitQuery(portfolioHistoryTableName, uuid, field, limit)
}
