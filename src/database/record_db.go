package database

import (
	"database/sql"
	"fmt"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects"
	"time"
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
	recordQuery             = `SELECT * from ` + recordTableName + ` ORDER BY time ASC;`

	recordDeleteRecord = `DELETE FROM ` + recordTableName + ` where uuid=$1`
)

func (d *Database) InitRecord() error {
	return d.Exec("record-init", recordTableCreateStatement)
}

func deleteRecord(record objects.Record, tx *sql.Tx) error {
	_, err := tx.Exec(recordDeleteRecord, record.Uuid)
	return err
}

func writeRecord(record objects.Record, tx *sql.Tx) error {
	_, err := tx.Exec(recordTableUpdateInsert, record.Uuid, record.SharePrice, record.RecordBookUuid, record.Fees, record.ShareCount, record.Taxes, record.Bonus, record.Result)
	return err
}

func (d *Database) GetRecords() ([]objects.Record, error) {
	var uuid, recordUuid string
	var sharePrice, fees, taxes, bonus, amount, id, result int64
	var t time.Time
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(recordQuery); err != nil {
		return nil, fmt.Errorf("failed to query portfolio err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	records := make([]objects.Record, 0)
	for rows.Next() {
		if err = rows.Scan(&id, &t, &uuid, &sharePrice, &recordUuid, &fees, &amount, &taxes, &bonus, &result); err != nil {
			return nil, err
		}
		records = append(records, objects.Record{
			Uuid: uuid,
			RecordBookUuid: recordUuid,
			ShareCount: amount,
			SharePrice: sharePrice,
			Taxes: taxes,
			Fees: fees,
			Bonus: bonus,
			Result:result,
			Time: t,
		})
	}
	return records, rows.Err()
}
