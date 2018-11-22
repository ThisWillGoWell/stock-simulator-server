package database

import (
	"encoding/json"
	"log"
	"time"

	"github.com/stock-simulator-server/src/effect"
)

var (
	effectTableName            = `effect`
	effectTableCreateStatement = `CREATE TABLE IF NOT EXISTS ` + effectTableName +
		`( ` +
		`uuid text NOT NULL, ` +
		`portfolio_uuid text NOT NULL, ` +
		`type text NOT NULL, ` +
		`title text NOT NULL, ` +
		`duration  bigint NOT NULL, ` +
		`start_time TIMESTAMPTZ NOT NULL, ` +
		`tag text NOT NULL, ` +
		`effect json NOT NULL, ` +
		`PRIMARY KEY(uuid)` +
		`);`

	effectTableUpdateInsert = `INSERT into ` + effectTableName + `(uuid, portfolio_uuid, type, title, duration, start_time, tag, effect) values($1, $2, $3, $4, $5, $6, $7, $8) ` +
		`ON CONFLICT (uuid) DO UPDATE SET effect=EXCLUDED.effect, title=EXCLUDED.title`

	effectTableQueryStatement  = "SELECT * FROM " + effectTableName + `;`
	effectTableDeleteStatement = "DELETE FROM " + effectTableName + " where uuid=$1"
)

func initEffect() {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		panic("could not begin effect init: " + err.Error())
	}
	_, err = tx.Exec(effectTableCreateStatement)
	if err != nil {
		tx.Rollback()
		panic("error occurred while creating effect table " + err.Error())
	}
	tx.Commit()
}

func writeEffect(entry *effect.Effect) {
	dbLock.Acquire("update-effect")
	defer dbLock.Release()
	tx, err := db.Begin()

	if err != nil {
		db.Close()
		panic("could not begin effect init" + err.Error())
	}
	e, err := json.Marshal(entry.InnerEffect)
	if err != nil {
	}
	_, err = tx.Exec(effectTableUpdateInsert, entry.Uuid, entry.PortfolioUuid, entry.Type, entry.Title, entry.Duration.Duration, entry.StartTime, entry.Tag, e)
	if err != nil {
		tx.Rollback()
		panic("error occurred while insert effect in table " + err.Error())
	}
	tx.Commit()
}

func populateEffects() {
	var effectType, effectJsonString, uuid, portfolioUuid, title, tag string
	var duration float64
	var startTime time.Time

	rows, err := db.Query(effectTableQueryStatement)
	if err != nil {
		log.Fatal("error reading effect database")
		panic("could not populate notifications: " + err.Error())
	}
	defer rows.Close()
	for rows.Next() {

		err := rows.Scan(&uuid, &portfolioUuid, &effectType, &title, &duration, &startTime, &tag, &effectJsonString)
		if err != nil {
			log.Fatal("error in querying ledger: ", err)
		}
		innerEffect := effect.UnmarshalJsonEffect(effectType, effectJsonString)
		effect.MakeEffect(uuid, portfolioUuid, title, effectType, tag, innerEffect, time.Duration(duration), startTime)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func deleteEffect(effect *effect.Effect) {
	tx, err := db.Begin()
	if err != nil {
		db.Close()
		panic("error opening db for deleting effect: " + err.Error())
	}
	_, err = tx.Exec(effectTableDeleteStatement, effect.Uuid)
	if err != nil {
		tx.Rollback()
		panic("error delete effect: " + err.Error())
	}
	tx.Commit()

}
