package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
)

func (d *Database) initItems() error {
	return d.Exec("init-items", itemsTableCreateStatement)
}

func (d *Database) WriteItem(entry *items.Item) error {

	innerItemStr, err := json.Marshal(entry.InnerItem)
	if err != nil {
		return fmt.Errorf("failed to marshal inner item err=[%v]", err)
	}
	return d.Exec("items-update", itemsTableUpdateInsert, entry.Uuid, entry.Type, entry.Name, entry.ConfigId, entry.PortfolioUuid, innerItemStr, entry.CreateTime)
}

func (d *Database) DeleteItem(item *items.Item) error {
	return d.Exec("items-delete", itemsTableDeleteStatement, item.Uuid)
}

func (d *Database) populateItems() error {
	var uuid, itemType, name, configId, portfolioUuid, innerJson string
	var createTime time.Time

	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(itemsTableQueryStatement); err != nil {
		return fmt.Errorf("failed to query portfolio err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		if err = rows.Scan(&uuid, &itemType, &name, &configId, &portfolioUuid, &innerJson, &createTime); err != nil {
			return err
		}
		if _, err = items.MakeItem(uuid, portfolioUuid, configId, itemType, name, innerJson, createTime); err != nil {
			return err
		}
	}
	return rows.Err()
}
