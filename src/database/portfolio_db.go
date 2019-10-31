package database

import (
	"database/sql"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/portfolio"
	"github.com/pkg/errors"
)

var (
	portfolioTableName            = `portfolio`
	portfolioTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + portfolioTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`name text NOT NULL,` +
		`wallet bigint NOT NULL,` +
		`level int NOT NULL, ` +
		`PRIMARY KEY(uuid) ` +
		`);`

	portfolioTableUpdateInsert = `INSERT into ` + portfolioTableName + `(uuid, name, wallet, level) values($1, $2, $3, $4) ` +
		`ON CONFLICT (uuid) DO UPDATE SET wallet=EXCLUDED.wallet,  level=EXCLUDED.level;`

	portfolioTableQueryStatement = "SELECT uuid, name, wallet, level FROM " + portfolioTableName + `;`

	portfolioHistoryTableName            = `portfolio_history`
	portfolioHistoryTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + portfolioHistoryTableName +
		`( ` +
		`time TIMESTAMPTZ NOT NULL,` +
		`uuid text NOT NULL,` +
		`net_worth bigint NULL,` +
		`wallet bigint NULL` +
		`);`

	portfolioHistoryTableUpdateInsert = `INSERT INTO ` + portfolioHistoryTableName + `(time, uuid, net_worth, wallet) values (NOW(),$1, $2, $3)`

	validPortfolioFields = map[string]bool{
		"wallet":    true,
		"net_worth": true,
	}
)

func (d *Database) InitPortfolio() error {
	if err := d.Exec("portfolio-init", portfolioTableCreateStatement); err != nil {
		return err
	}
	return d.Exec("portfolio-history-init", portfolioHistoryTableCreateStatement)
}

func (d *Database) WritePortfolio(port *portfolio.Portfolio) error {
	if err := d.Exec(portfolioTableUpdateInsert, port.Uuid, port.UserUUID, port.Wallet, port.Level); err != nil {
		return err
	}
	return d.Exec(portfolioHistoryTableUpdateInsert, port.Uuid, port.NetWorth, port.Wallet)
}

func (d *Database) populatePortfolios() error {
	var uuid, userUuid string
	var wallet, level int64
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(portfolioTableQueryStatement); err != nil {
		return fmt.Errorf("failed to query portfolio err=%v", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		if err = rows.Scan(&uuid, &userUuid, &wallet, &level); err != nil {
			return err
		}
		if _, err = portfolio.MakePortfolio(uuid, userUuid, wallet, level); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (d *Database) MakePortfolioHistoryTimeQuery(uuid, timeLength, field, intervalLength string) ([][]interface{}, error) {
	if _, valid := validPortfolioFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return MakeHistoryTimeQuery(portfolioHistoryTableName, uuid, timeLength, field, intervalLength)

}

func (d *Database) MakePortfolioHistoryLimitQuery(uuid, field string, limit int) ([][]interface{}, error) {
	if _, valid := validPortfolioFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return MakeHistoryLimitQuery(portfolioHistoryTableName, uuid, field, limit)
}
