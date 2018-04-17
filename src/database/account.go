package database

import (
	"github.com/stock-simulator-server/src/account"
	"log"
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
		`PRIMARY KEY(uuid)` +
		`);`

	accountTableUpdateInsert = `INSERT into ` + accountTableName + `(uuid, name, display_name, password, portfolio_uuid) values($1, $2, $3, $4, $5) ` +
		`ON CONFLICT (uuid) DO UPDATE SET display_name=EXCLUDED.display_name, password=EXCLUDED.password;`

	accountTableQueryStatement = "SELECT uuid, name, display_name, password, portfolio_uuid FROM " + accountTableName + `;`
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
		panic("could not begin stocks init")
	}
	_, err = tx.Exec(accountTableUpdateInsert, user.Uuid, user.UserName, user.DisplayName, user.Password, user.PortfolioId)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert account in table " + err.Error())
	}
	tx.Commit()
}

func populateUsers() {
	var uuid, name, displayName, password, portfolioId string

	rows, err := db.Query(accountTableQueryStatement)
	if err != nil {
		log.Fatal("error quiering databse")
		panic("could not populate users: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &name, &displayName, &password, &portfolioId)
		if err != nil {
			log.Fatal(err)
		}
		account.MakeUser(uuid, name, displayName, password, portfolioId)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
