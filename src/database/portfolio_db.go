package database

import (
	"database/sql"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"
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

	portfolioTableDelete        = "DELETE from " + portfolioTableName + `WHERE uuid = $1`
	portfolioHistroyTableDelete = "DELETE from " + portfolioHistoryTableName + `WHERE uuid = $1`

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

func writePortfolio(port models.Portfolio, tx *sql.Tx) error {
	if _, err := tx.Exec(portfolioTableUpdateInsert, port.Uuid, port.UserUUID, port.Wallet, port.Level); err != nil {
		return err
	}
	if _, err := tx.Exec(portfolioHistoryTableUpdateInsert, port.Uuid, port.NetWorth, port.Wallet); err != nil {
		return err
	}
	return nil
}

func deletePortfolio(port models.Portfolio, tx *sql.Tx) error {
	if _, err := tx.Exec(portfolioTableDelete, port.Uuid); err != nil {
		return err
	}
	if _, err := tx.Exec(portfolioHistroyTableDelete, port.Uuid); err != nil {
		return nil
	}
	return nil
}

func (d *Database) PopulatePortfolios() (map[string]models.Portfolio, error) {
	var uuid, userUuid string
	var wallet, level int64
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(portfolioTableQueryStatement); err != nil {
		return nil, fmt.Errorf("failed to query portfolio err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	ports := make(map[string]models.Portfolio)
	for rows.Next() {
		if err = rows.Scan(&uuid, &userUuid, &wallet, &level); err != nil {
			return nil, err
		}
		ports[uuid] = models.Portfolio{
			UserUUID: userUuid,
			Uuid:     uuid,
			Wallet:   wallet,
			NetWorth: 0,
			Level:    level,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ports, nil
}

func (d *Database) MakePortfolioHistoryTimeQuery(uuid, timeLength, field, intervalLength string) ([][]interface{}, error) {
	if _, valid := validPortfolioFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return d.MakeHistoryTimeQuery(portfolioHistoryTableName, uuid, timeLength, field, intervalLength)

}

func (d *Database) MakePortfolioHistoryLimitQuery(uuid, field string, limit int) ([][]interface{}, error) {
	if _, valid := validPortfolioFields[field]; !valid {
		return nil, errors.New("not valid choice")
	}
	return d.MakeHistoryLimitQuery(portfolioHistoryTableName, uuid, field, limit)
}
