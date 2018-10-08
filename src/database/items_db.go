package database

import (
	"encoding/json"
	"log"

	"github.com/stock-simulator-server/src/items"
)

var (
	itemsTableName            = `items`
	itemsTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + itemsTableName +
		`( ` +
		`uuid text NOT NULL, ` +
		`type text NOT NULL, ` +
		`item json NOT NULL, ` +
		`PRIMARY KEY(uuid)` +
		`);`

	itemsTableUpdateInsert = `INSERT into ` + itemsTableName + `(uuid, type, item) values($1, $2, $3) ` +
		`ON CONFLICT (uuid) DO UPDATE SET item=EXCLUDED.item`

	itemsTableQueryStatement = "SELECT type, item FROM " + itemsTableName + `;`
	//getCurrentPrice()
)

func initItems() {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		panic("could not begin ledger init: " + err.Error())
	}
	_, err = tx.Exec(itemsTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating leger table " + err.Error())
	}
	tx.Commit()
}

func writeItem(entry items.Item) {
	dbLock.Acquire("update-item")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin item init" + err.Error())
	}
	item, err := json.Marshal(entry)
	if err != nil {
	}

	_, err = tx.Exec(itemsTableUpdateInsert, entry.GetUuid(), entry.GetType(), item)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert item in table " + err.Error())
	}
	tx.Commit()
}

func populateItems() {
	var itemType, itemJsonString string

	rows, err := db.Query(itemsTableQueryStatement)
	if err != nil {
		log.Fatal("error reading items databse")
		panic("could not populate notifications: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&itemType, &itemJsonString)
		if err != nil {
			log.Fatal("error in querying ledger: ", err)
		}
		item := items.UnmarshalJsonItem(itemType, itemJsonString)
		items.LoadItem(item)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
