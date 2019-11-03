package database

import (
	"database/sql"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"

	"github.com/ThisWillGoWell/stock-simulator-server/src/user"
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

func (d *Database) WriteUser(user models.User) error {
	return d.Exec("user-update", accountTableUpdateInsert, user.Uuid, user.UserName, user.DisplayName, user.Password, user.PortfolioId, user.ConfigStr)
}

func (d *Database) DeleteUser(uuid string) error {
	return d.Exec("user-delete", userTableDelete, uuid)
}

func (d *Database) populateUsers() error {
	var uuid, name, displayName, password, portfolioId, config string
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(effectTableQueryStatement); err != nil {
		return fmt.Errorf("failed to query portfolio err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		if err = rows.Scan(&uuid, &name, &displayName, &password, &portfolioId, &config); err != nil {
			return err
		}
		if _, err = user.MakeUser(uuid, name, displayName, password, portfolioId, config); err != nil {
			return err
		}
	}
	return rows.Err()
}
