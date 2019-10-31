package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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

func (d *Database) WriteNotification(entry *notification.Notification) error {
	jsonString, err := json.Marshal(entry.Notification)
	if err != nil {
		return fmt.Errorf("failed to marshal inner notificaion err=%v", err)
	}

	return d.Exec(notificationTableUpdateInsert, entry.Uuid, entry.PortfolioUuid, entry.Seen, entry.Type, entry.Timestamp, jsonString)
}

func (d *Database) DeleteNotification(note *notification.Notification) error {
	return d.Exec("notification-delete", notificationDeleteStatement, note.Uuid)
}

func (d *Database) populateNotification() {
	var uuid, userUuid, jsonString, notType string
	var seen bool
	var t time.Time

	rows, err := d.db.Query(notificationTableQueryStatement)
	if err != nil {
		log.Fatal("error reading notifications databse")
		panic("could not populate notifications: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &userUuid, &seen, &t, &notType, &jsonString)
		if err != nil {
			log.Fatal("error in querying ledger: ", err)
		}
		note := notification.JsonToNotification(jsonString, notType)
		notification.MakeNotification(uuid, userUuid, notType, t, seen, note)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
