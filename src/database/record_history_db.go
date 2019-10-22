package database

import (
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/log"

	"github.com/ThisWillGoWell/stock-simulator-server/src/record"
)

var (
	recordHistoryTableName            = `record_history`
	recordHistoryTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + recordHistoryTableName +
		`( ` +
		`id serial, ` +
		`time TIMESTAMPTZ NOT NULL,` +
		`uuid text NOT NULL,` +
		`share_price bigint NOT NULL,` +
		`record_uuid text NOT NULL, ` +
		`fees bigint NULL, ` +
		`amount bigint NULL,` +
		`taxes bigint NULL, ` +
		`bonus bigint NULL, ` +
		`result bigint NULL` +
		`);`

	recordHistoryTableUpdateInsert = `INSERT INTO ` + recordHistoryTableName + `(time, uuid, share_price, record_uuid, fees, amount, taxes, bonus, result) values (NOW(), $1, $2, $3, $4, $5, $6, $7, $8);`
	recordHistoryQuery             = `SELECT * from ` + recordHistoryTableName
	//getCurrentPrice()
)

func initRecordHistory() {
	tx, err := db.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin record history init: " + err.Error())
	}
	_, err = tx.Exec(recordHistoryTableCreateStatement)
	if err != nil {
		log.Log.Printf("during create record history err=%v\n", err)
	}
}

func writeRecordHistory(record *record.Record) {
	tx, err := db.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin record history init: " + err.Error())
	}
	_, err = tx.Exec(recordHistoryTableUpdateInsert, record.Uuid, record.SharePrice, record.RecordBookUuid, record.Fees, record.ShareCount, record.Taxes, record.Bonus, record.Result)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert record in table " + err.Error())
	}
	tx.Commit()
}

func populateRecords() {
	var uuid, recordUuid string
	var sharePrice, fees, taxes, bonus, amount, id, result int64
	var t time.Time
	rows, err := db.Query(recordHistoryQuery)
	if err != nil {
		log.Log.Fatal("error query records database", err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &t, &uuid, &sharePrice, &recordUuid, &fees, &amount, &taxes, &bonus, &result)
		if err != nil {
			log.Log.Fatal(err)
		}

		record.MakeRecord(uuid, recordUuid, amount, sharePrice, taxes, fees, bonus, result, t)
	}
	err = rows.Err()
	if err != nil {
		log.Log.Fatal(err)
	}
}
