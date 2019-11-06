package database

import (
	"database/sql"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"
)

var (
	accountTableName            = `user`
	accountTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + accountTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`name text NOT NULL,` +
		`display_name text NOT NULL,` +
		`password text NOT NULL,` +
		`portfolio_uuid text NOT NULL,` +
		`config json NULL, ` +
		`PRIMARY KEY(uuid)` +
		`);`

	accountTableUpdateInsert = `INSERT into ` + accountTableName + `(uuid, name, display_name, password, portfolio_uuid, config) values($1, $2, $3, $4, $5, $6) ` +
		`ON CONFLICT (uuid) DO UPDATE SET display_name=EXCLUDED.display_name, password=EXCLUDED.password, config=EXCLUDED.config;`

	accountTableQueryStatement = "SELECT uuid, name, display_name, password, portfolio_uuid, config FROM " + accountTableName + `;`
	userTableDelete            = "DELETE from " + accountTableName + `WHERE uuid = $1`
)

func (d *Database) initAccount() error {
	return d.Exec("user-init", accountTableCreateStatement)
}

func writeUser(user models.User, tx *sql.Tx) error {
	_, err := tx.Exec(accountTableUpdateInsert, user.Uuid, user.UserName, user.DisplayName, user.Password, user.PortfolioId, user.ConfigStr)
	return err
}

func deleteUser(user models.User, tx *sql.Tx) error {
	_, err := tx.Exec(userTableDelete, uuid)
	return err
}

func (d *Database) PopulateUsers() (map[string]models.User, error) {
	var uuid, name, displayName, password, portfolioId, config string
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(effectTableQueryStatement); err != nil {
		return nil, fmt.Errorf("failed to query portfolio err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	users := make(map[string]models.User)

	for rows.Next() {
		if err = rows.Scan(&uuid, &name, &displayName, &password, &portfolioId, &config); err != nil {
			return nil, err
		}
		users[uuid] = models.User{
			UserName:      name,
			Password:      password,
			DisplayName:   displayName,
			Uuid:          uuid,
			PortfolioId:   portfolioId,
			Active:        false,
			Config:        nil,
			ConfigStr:     config,
			ActiveClients: 0,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, err
}
