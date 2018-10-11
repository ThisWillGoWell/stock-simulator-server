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

var queryChan = make(chan interface{})

type historicalTimeQuery struct {
	table          string
	uuid           string
	timeLength     string
	intervalLength string
	field          string
	responseChan   chan *historicalQueryResponse
}

type historicalLimitQuery struct {
	table        string
	uuid         string
	field        string
	limit        int
	responseChan chan *historicalQueryResponse
}

type historicalQueryResponse struct {
	response [][]interface{}
	err      error
}

func runHistoricalQueries() {
	go func() {
		for {
			query := <-queryChan
			switch query.(type) {
			case *historicalTimeQuery:
				response, err := makeHistoryTimeQuery(query.(*historicalTimeQuery))
				query.(*historicalTimeQuery).responseChan <- &historicalQueryResponse{response: response, err: err}
			case *historicalLimitQuery:
				response, err := makeHistoricalLimitQuery(query.(*historicalLimitQuery))
				query.(*historicalLimitQuery).responseChan <- &historicalQueryResponse{response: response, err: err}
			}
		}
	}()
}

func MakeHistoryTimeQuery(table, uuid, timeLength, field, intervalLength string) ([][]interface{}, error) {
	q := &historicalTimeQuery{table, uuid, timeLength, intervalLength, field, make(chan *historicalQueryResponse, 1)}
	queryChan <- q
	r := <-q.responseChan
	return r.response, r.err
}
func makeHistoryTimeQuery(query *historicalTimeQuery) ([][]interface{}, error) {
	//rows, err := tx.Query("SELECT time_bucket('60 seconds', time) AS tb, AVG(current_price) AS val FROM stocks_history  WHERE time > NOW() - interval '600 seconds' and uuid='E30B70AD77B26C' GROUP BY tb  ORDER BY tb DESC")
	querySmt := fmt.Sprintf(historyTableTimeQuery, query.intervalLength, query.field, query.table, query.timeLength)
	rows, err := db.Query(querySmt, query.uuid)
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
func MakeHistoryLimitQuery(table, uuid, field string, limit int) ([][]interface{}, error) {
	q := &historicalLimitQuery{table, uuid, field, limit, make(chan *historicalQueryResponse)}
	queryChan <- q
	r := <-q.responseChan
	return r.response, r.err
}

func makeHistoricalLimitQuery(query *historicalLimitQuery) ([][]interface{}, error) {
	querySmt := fmt.Sprintf(historyTableLimitQuery, query.field, query.table)
	rows, err := db.Query(querySmt, query.uuid, query.limit)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	return rowsToResponse(rows)
}
