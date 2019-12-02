package database

import (
	"database/sql"
	"fmt"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects"
)

var (
	accountTableName            = `users`
	accountTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + accountTableName +
		`( ` +
		`id serial,` +
		`uuid text NOT NULL,` +
		`name text NOT NULL,` +
		`display_name text NOT NULL,` +
		`password text NOT NULL,` +
		`portfolio_uuid text NOT NULL,` +
		`config json NULL, ` +
		`is_admin boolean NOT NULL DEFAULT FALSE` +
		`PRIMARY KEY(uuid)` +
		`);`

	accountTableUpdateInsert = `INSERT into ` + accountTableName + `(uuid, name, display_name, password, portfolio_uuid, config, isAdmin) values($1, $2, $3, $4, $5, $6, $7) ` +
		`ON CONFLICT (uuid) DO UPDATE SET display_name=EXCLUDED.display_name, password=EXCLUDED.password, config=EXCLUDED.config, is_admin=EXCLUDED.isAdmin;`

	userTableQueryStatement = "SELECT uuid, name, display_name, password, portfolio_uuid, config, is_admin FROM " + accountTableName + `;`
	userTableDelete         = "DELETE from " + accountTableName + `WHERE uuid = $1`
)

func (d *Database) initAccount() error {
	return d.Exec("user-init", accountTableCreateStatement)
}

func writeUser(user objects.User, tx *sql.Tx) error {
	_, err := tx.Exec(accountTableUpdateInsert, user.Uuid, user.UserName, user.DisplayName, user.Password, user.PortfolioId, user.ConfigStr, user.IsAdmin)
	return err
}

func deleteUser(user objects.User, tx *sql.Tx) error {
	_, err := tx.Exec(userTableDelete, user.Uuid)
	return err
}

func (d *Database) GetUsers() ([]objects.User, error) {
	var uuid, name, displayName, password, portfolioId, config string
	var isAdmin bool
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(userTableQueryStatement); err != nil {
		return nil, fmt.Errorf("failed to query portfolio err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	users := make([]objects.User, 0)

	for rows.Next() {
		if err = rows.Scan(&uuid, &name, &displayName, &password, &portfolioId, &config, &isAdmin); err != nil {
			return nil, err
		}
		users = append(users, objects.User{
			UserName:      name,
			Password:      password,
			DisplayName:   displayName,
			Uuid:          uuid,
			PortfolioId:   portfolioId,
			Active:        false,
			Config:        nil,
			ConfigStr:     config,
			IsAdmin:       isAdmin,
			ActiveClients: 0,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, err
}
