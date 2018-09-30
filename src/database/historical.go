package database

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	historyTableTimeQuery  = `SELECT time_bucket('%s', time) AS tb, AVG(%s) AS val FROM %s WHERE time > NOW() - interval '%s' and uuid=$1 GROUP BY tb  ORDER BY tb DESC`
	historyTableLimitQuery = `SELECT time, %s FROM %s WHERE uuid=$1 LIMIT $2;`
)

func MakeHistoryTimeQuery(table, uuid, timeLength, field, intervalLength string) ([][]interface{}, error) {
	tx, err := ts.Begin()
	if err != nil {
		return nil, err
	}

	//rows, err := tx.Query("SELECT time_bucket('60 seconds', time) AS tb, AVG(current_price) AS val FROM stocks_history  WHERE time > NOW() - interval '600 seconds' and uuid='E30B70AD77B26C' GROUP BY tb  ORDER BY tb DESC")
	querySmt := fmt.Sprintf(historyTableTimeQuery, intervalLength, field, table, timeLength)
	rows, err := tx.Query(querySmt, uuid)
	if err != nil {
		return nil, err
	}
	return rowsToResponse(rows)

}

func rowsToResponse(rows *sql.Rows) ([][]interface{}, error) {
	var t time.Time
	var value float64
	response := make([][]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(&t, &value)
		if err != nil {
			return nil, err
		}
		response = append(response, []interface{}{t, value})
	}
	return response, nil
}
func MakeHistoryLimitQuery(tableName, uuid, field string, limit int) ([][]interface{}, error) {

	querySmt := fmt.Sprintf(historyTableLimitQuery, field, tableName)
	rows, err := ts.Query(querySmt, uuid, limit)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return rowsToResponse(rows)
}
