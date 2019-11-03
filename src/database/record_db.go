package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/record"
)

var (
	recordTableName            = `records`
	recordTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + recordTableName +
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

	recordTableUpdateInsert = `INSERT INTO ` + recordTableName + `(time, uuid, share_price, record_uuid, fees, amount, taxes, bonus, result) values (NOW(), $1, $2, $3, $4, $5, $6, $7, $8);`
	recordQuery             = `SELECT * from ` + recordTableName

	recordDeleteRecord = `DELETE FROM ` + recordTableName + ` where uuid=$1`
)

func (d *Database) InitRecord() error {
	return d.Exec("record-init", recordTableCreateStatement)
}

func (d *Database) DeleteRecord(record *record.Record) error {
	return d.Exec("record-delete", recordDeleteRecord, record.Uuid)
}

func (d *Database) WriteRecord(record *record.Record) error {
	return d.Exec("record-update", recordTableUpdateInsert, record.Uuid, record.SharePrice, record.RecordBookUuid, record.Fees, record.ShareCount, record.Taxes, record.Bonus, record.Result)
}

func (d *Database) populateRecords() error {
	var uuid, recordUuid string
	var sharePrice, fees, taxes, bonus, amount, id, result int64
	var t time.Time
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(recordQuery); err != nil {
		return fmt.Errorf("failed to query portfolio err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		if err = rows.Scan(&id, &t, &uuid, &sharePrice, &recordUuid, &fees, &amount, &taxes, &bonus, &result); err != nil {
			return err
		}
		if _, err = record.MakeRecord(uuid, recordUuid, amount, sharePrice, taxes, fees, bonus, result, t); err != nil {
			return err
		}
	}
	return rows.Err()
}
