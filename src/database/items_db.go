package database

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/items"
)

var (
	itemsTableName            = `items`
	itemsTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + itemsTableName +
		`( ` +
		`uuid text NOT NULL, ` +
		`type text NOT NULL, ` +
		`name text NOT NULL, ` +
		`config_id text NOT NULL, ` +
		`portfolio_uuid text NOT NULL, ` +
		`inner_item json NOT NULL, ` +
		`create_time TIMESTAMPTZ NOT NULL, ` +
		`PRIMARY KEY(uuid)` +
		`);`

	itemsTableUpdateInsert = `INSERT into ` + itemsTableName + `(uuid, type, name, config_id, portfolio_uuid, inner_item,create_time) values($1, $2, $3, $4, $5, $6, $7) ` +
		`ON CONFLICT (uuid) DO UPDATE SET inner_item=EXCLUDED.inner_item`

	itemsTableQueryStatement  = "SELECT * FROM " + itemsTableName + `;`
	itemsTableDeleteStatement = "DELETE FROM " + itemsTableName + " where uuid=$1"
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

func writeItem(entry *items.Item) {
	dbLock.Acquire("update-item")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin item init" + err.Error())
	}
	innerItemStr, err := json.Marshal(entry.InnerItem)
	if err != nil {
	}

	_, err = tx.Exec(itemsTableUpdateInsert, entry.Uuid, entry.Type, entry.Name, entry.ConfigId, entry.PortfolioUuid, innerItemStr, entry.CreateTime)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert item in table " + err.Error())
	}
	tx.Commit()
}

func populateItems() {
	var uuid, itemType, name, configId, portfolioUuid, innerJson string
	var createTime time.Time
	rows, err := db.Query(itemsTableQueryStatement)
	if err != nil {
		log.Fatal("error reading items database")
		panic("could not populate notifications: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &itemType, &name, &configId, &portfolioUuid, &innerJson, &createTime)
		if err != nil {
			log.Fatal("error in querying ledger: ", err)
		}
		items.MakeItem(uuid, portfolioUuid, configId, itemType, name, innerJson, createTime)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func deleteItem(item *items.Item) {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		panic("error opening db for deleting item: " + err.Error())
	}
	_, err = tx.Exec(itemsTableDeleteStatement, item.Uuid)
	if err != nil {
		tx.Rollback()
		panic("error delete item: " + err.Error())
	}
	tx.Commit()

}
