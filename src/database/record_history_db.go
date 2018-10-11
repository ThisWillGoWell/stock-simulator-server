package database

import (
	"log"
	"time"

	"github.com/stock-simulator-server/src/record"
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
		`bonus bigint NULL ` +
		`);`
	recordHistoryTSInit = `CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE; SELECT create_hypertable('` + recordHistoryTableName + `', 'time');`

	recordHistoryTableUpdateInsert = `INSERT INTO ` + recordHistoryTableName + `(time, uuid, share_price, record_uuid, fees, amount, taxes, bonus) values (NOW(), $1, $2, $3, $4, $5, $6, $7);`
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

	}
	tx.Commit()
	tx, err = db.Begin()
	_, err = tx.Exec(recordHistoryTSInit)
	if err != nil {

	}
	tx.Commit()
}

func writeRecordHistory(record *record.Record) {
	tx, err := db.Begin()
	if err != nil {
		ts.Close()
		panic("could not begin record history init: " + err.Error())
	}
	_, err = tx.Exec(recordHistoryTableUpdateInsert, record.Uuid, record.SharePrice, record.RecordUuid, record.Fees, record.Amount, record.Taxes, record.Bonus)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert record in table " + err.Error())
	}
	tx.Commit()
}

func populateRecords() {
	var uuid, recordUuid string
	var sharePrice, fees, taxes, bonus, amount, id int64
	var t time.Time
	rows, err := db.Query(recordHistoryQuery)
	if err != nil {
		log.Fatal("error query records database", err)
		panic("could not populate records: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &t, &uuid, &sharePrice, &recordUuid, &fees, &amount, &taxes, &bonus)
		if err != nil {
			panic(err)
			log.Fatal(err)
		}
		record.MakeRecord(uuid, recordUuid, amount, sharePrice, taxes, fees, bonus, t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
