package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"
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

func writeItem(entry models.Item, tx *sql.Tx) error {
	innerItemStr, err := json.Marshal(entry.InnerItem)
	if err != nil {
		return fmt.Errorf("failed to marshal inner item err=[%v]", err)
	}
	_, err = tx.Exec(itemsTableUpdateInsert, entry.Uuid, entry.Type, entry.Name, entry.ConfigId, entry.PortfolioUuid, innerItemStr, entry.CreateTime)
	return err
}

func deleteItem(entry models.Item, tx *sql.Tx) error {
	_, err := tx.Exec(itemsTableDeleteStatement, entry.Uuid)
	return err
}

func (d *Database) GetItems() (map[string]models.Item, error) {
	var uuid, itemType, name, configId, portfolioUuid, innerJson string
	var createTime time.Time

	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(itemsTableQueryStatement); err != nil {
		return nil, fmt.Errorf("failed to query items err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	items := make(map[string]models.Item)
	for rows.Next() {
		if err = rows.Scan(&uuid, &itemType, &name, &configId, &portfolioUuid, &innerJson, &createTime); err != nil {
			return nil, err
		}
		items[uuid] = models.Item{
			Uuid:          uuid,
			Name:          name,
			ConfigId:      configId,
			Type:          itemType,
			PortfolioUuid: portfolioUuid,
			CreateTime:    createTime,
			InnerItem:     innerJson,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
