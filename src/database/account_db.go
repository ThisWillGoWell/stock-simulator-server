package database

import (
	"log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/account"
)

var (
	accountTableName            = `account`
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
	//todo seperate inset for json, not each time
	accountTableUpdateInsert = `INSERT into ` + accountTableName + `(uuid, name, display_name, password, portfolio_uuid, config) values($1, $2, $3, $4, $5, $6) ` +
		`ON CONFLICT (uuid) DO UPDATE SET display_name=EXCLUDED.display_name, password=EXCLUDED.password, config=EXCLUDED.config;`

	accountTableQueryStatement = "SELECT uuid, name, display_name, password, portfolio_uuid, config FROM " + accountTableName + `;`
	//getCurrentPrice()
)

func initAccount() {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		panic("could not begin account db init: " + err.Error())
	}
	_, err = tx.Exec(accountTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating account table " + err.Error())
	}
	tx.Commit()
}

func writeUser(user *account.User) {
	dbLock.Acquire("add user")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin user init: " + err.Error())
	}

	_, err = tx.Exec(accountTableUpdateInsert, user.Uuid, user.UserName, user.DisplayName, user.Password, user.PortfolioId, user.ConfigStr)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert account in table " + err.Error())
	}
	tx.Commit()
}

func populateUsers() {
	var uuid, name, displayName, password, portfolioId, config string

	rows, err := db.Query(accountTableQueryStatement)
	if err != nil {
		log.Fatal("error quiering databse", err)
		panic("could not populate users: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &name, &displayName, &password, &portfolioId, &config)
		if err != nil {
			log.Fatal("error :", err)
		}
		account.MakeUser(uuid, name, displayName, password, portfolioId, config)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
