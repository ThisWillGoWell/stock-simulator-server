package database

import (
	"encoding/json"
	"log"
	"time"

	"github.com/stock-simulator-server/src/notification"
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

	notificationTableUpdateInsert = `INSERT into ` + notificationTableName + `(uuid, userUuid, seen, timestamp, type, notification) values($1, $2, $3, $4, $5) ` +
		`ON CONFLICT (uuid) DO UPDATE SET seen=EXCLUDED.seen`

	notificationTableQueryStatement = "SELECT uuid, portfolio_id, stock_id, amount, investment_value FROM " + notificationTableName + `;`
	//getCurrentPrice()
)

func initNotification() {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		panic("could not begin ledger init: " + err.Error())
	}
	_, err = tx.Exec(notificationTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating leger table " + err.Error())
	}
	tx.Commit()
}

func writeNotification(entry *notification.Notification) {
	dbLock.Acquire("update-notification")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin notification init" + err.Error())
	}
	jsonString, err := json.Marshal(entry.Notification)
	if err != nil {
	}

	_, err = tx.Exec(notificationTableUpdateInsert, entry.Uuid, entry.UserUuid, entry.Seen, entry.Type, entry.Timestamp, jsonString)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert notification in table " + err.Error())
	}
	tx.Commit()
}

func populateNotification() {
	var uuid, userUuid, jsonString, notType string
	var seen bool
	var t time.Time

	rows, err := db.Query(notificationTableQueryStatement)
	if err != nil {
		log.Fatal("error databse")
		panic("could not populate notifications: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&uuid, &userUuid, &seen, &t, &notType, &jsonString)
		if err != nil {
			log.Fatal("error in querying ledger: ", err)
		}
		note := notification.JsonToNotifcation(jsonString, notType)
		notification.MakeNotification(uuid, userUuid, notType, t, seen, note)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
