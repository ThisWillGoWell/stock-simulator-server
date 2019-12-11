package database

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	historyTableTimeQuery  = `SELECT time AS tb, AVG(%s) AS current_price FROM %s WHERE time > NOW() - interval '%s' and uuid=$1 GROUP BY tb  ORDER BY tb DESC`
	historyTableLimitQuery = `SELECT time, %s FROM %s WHERE uuid=$1 LIMIT $2;`
)

func (d *Database) MakeHistoryTimeQuery(table, uuid, timeLength, field, intervalLength string) ([][]interface{}, error) {
	querySmt := fmt.Sprintf(historyTableTimeQuery, field, table, timeLength)
	rows, err := d.db.Query(querySmt, uuid)
	if err != nil {
		return nil, err
	}
	return rowsToResponse(rows)
}

func (d *Database) MakeHistoryLimitQuery(table, uuid, field string, limit int) ([][]interface{}, error) {
	querySmt := fmt.Sprintf(historyTableLimitQuery, field, table)
	var rows *sql.Rows
	var err error
	if rows, err = d.db.Query(querySmt, uuid, limit); err != nil {
		return nil, err
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}
	return rowsToResponse(rows)
}

func rowsToResponse(rows *sql.Rows) ([][]interface{}, error) {
	var t time.Time
	var value float64
	var err error
	response := make([][]interface{}, 0)
	for rows.Next() {
		if err = rows.Scan(&t, &value); err != nil {
			return nil, err
		}
		response = append(response, []interface{}{t, value})
	}
	return response, nil
}
