package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/utils"

	"github.com/ThisWillGoWell/stock-simulator-server/src/models"
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

func (d *Database) initEffect() error {
	return d.Exec("init-effect", effectTableCreateStatement)
}

func writeEffect(entry models.Effect, tx *sql.Tx) error {
	e, err := json.Marshal(entry.InnerEffect)
	if err != nil {
		return fmt.Errorf("failed to marshal inner effect err=[%v]", err)
	}
	_, err = tx.Exec(effectTableUpdateInsert, entry.Uuid, entry.PortfolioUuid, entry.Type, entry.Title, entry.Duration.Duration, entry.StartTime, entry.Tag, e)
	return err
}

func deleteEffect(entry models.Effect, tx *sql.Tx) error {
	_, err := tx.Exec(effectTableDeleteStatement, entry.Uuid)
	return err
}

func (d *Database) GetEffects() (map[string]models.Effect, error) {
	var effectType, effectJsonString, uuid, portfolioUuid, title, tag string
	var duration float64
	var startTime time.Time
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(effectTableQueryStatement); err != nil {
		return nil, fmt.Errorf("failed to query portfolio err=[%v]", err)
	}
	defer func() {
		_ = rows.Close()
	}()
	effects := make(map[string]models.Effect)

	for rows.Next() {
		if err = rows.Scan(&uuid, &portfolioUuid, &effectType, &title, &duration, &startTime, &tag, &effectJsonString); err != nil {
			return nil, err
		}
		effects[uuid] = models.Effect{
			PortfolioUuid: portfolioUuid,
			Uuid:          uuid,
			Title:         title,
			Duration:      utils.Duration{Duration: time.Duration(duration)},
			StartTime:     startTime,
			Type:          title,
			InnerEffect:   effectJsonString,
			Tag:           tag,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return effects, nil

}
