package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"

	"github.com/ThisWillGoWell/stock-simulator-server/src/notification"
)

var (
	notificationTableName            = `notification`
	notificationTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + notificationTableName +
		`( ` +
		`uuid text NOT NULL, ` +
		`userUuid text NOT NULL, ` +
		`seen bool NOT NULL, ` +
		`timestamp TIMESTAMPTZ NOT NULL, ` +
		`type text NOT NULL, ` +
		`notification json NOT NULL, ` +
		`PRIMARY KEY(uuid)` +
		`);`

	notificationTableUpdateInsert = `INSERT into ` + notificationTableName + `(uuid, userUuid, seen, type, timestamp, notification) values($1, $2, $3, $4, $5, $6) ` +
		`ON CONFLICT (uuid) DO UPDATE SET seen=EXCLUDED.seen`

	notificationTableQueryStatement = "SELECT uuid, userUuid, seen, timestamp, type, notification FROM " + notificationTableName + `;`

	notificationDeleteStatement = `DELETE from ` + notificationTableName + ` where uuid=$1`
)

func (d *Database) InitNotification() error {
	return d.Exec("notification-init", notificationTableCreateStatement)
}

func writeNotification(entry models.Notification, tx *sql.Tx) error {
	jsonString, err := json.Marshal(entry.Notification)
	if err != nil {
		return fmt.Errorf("failed to marshal inner notificaion err=[%v]", err)
	}
	_, err = tx.Exec(notificationTableUpdateInsert, entry.Uuid, entry.PortfolioUuid, entry.Seen, entry.Type, entry.Timestamp, jsonString)
	return err
}

func deleteNotification(note models.Notification, tx *sql.Tx) error {
	_, err := tx.Exec(notificationDeleteStatement, note.Uuid)
	return err
}

func (d *Database) GetNotification() (map[string]models.Notification, error ) {
	var uuid, userUuid, jsonString, notType string
	var seen bool
	var t time.Time

	rows, err := d.db.Query(notificationTableQueryStatement)
	if err != nil {
		return nil, fmt.Errorf("failed to query database err=[%v]", err)
	}
	n := make(map[string]models.Notification)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &userUuid, &seen, &t, &notType, &jsonString)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notificaion err=[%v]", err )
		}
		n[uuid] = models.Notification{
			Uuid: uuid,
			PortfolioUuid: userUuid,
			Timestamp: t,
			Seen: seen,
			Type: notType,
			Notification: jsonString,
		}
	}
	return n, rows.Err()
}
