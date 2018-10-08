package histroy

import (
	"errors"
	"fmt"
	"time"

	"github.com/stock-simulator-server/src/lock"

	"github.com/stock-simulator-server/src/database"
	"github.com/stock-simulator-server/src/ledger"
	"github.com/stock-simulator-server/src/messages"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/valuable"
)

type queryCacheItem struct {
	query          *Query
	response       *messages.QueryResponse
	lastUpdateTime time.Time
	lastUseTime    time.Time
	validTime      time.Duration
}

var queryCache = make(map[string]*queryCacheItem)
var queryCacheLock = lock.NewLock("query-cache-lock")
var expirationTime = time.Hour * 48

func RunCacheUpdater() {
	go func() {
		for {
			queryCacheLock.Acquire("clean")
			for queryKey, queryItem := range queryCache {
				if time.Since(queryItem.lastUseTime) > expirationTime {
					delete(queryCache, queryKey)
				} else {
					if time.Since(queryItem.lastUpdateTime) > queryItem.validTime {
						fmt.Println("updating query:", queryKey, time.Now())
						go makeQuery(queryItem.query)
						go func() {
							<-queryItem.query.ResponseChannel
						}()
					}
				}
			}
			queryCacheLock.Release()
			<-time.After(time.Minute * 1)
		}
	}()
}

func makeQueryHash(q *Query) string {
	return fmt.Sprintf("%s-%s-%s-%d-%s-%s", q.Type, q.QueryUUID, q.QueryField, q.Limit, q.TimeLength, q.TimeInterval)
}

type Query struct {
	Message         *messages.QueryMessage
	Type            string
	QueryUUID       string
	QueryField      string
	TimeInterval    string
	Interval        time.Duration
	TimeLength      string
	Limit           int
	ResponseChannel chan *messages.QueryResponse
}

func BuildQuery(qm *messages.QueryMessage) *Query {
	duration := qm.Length.Duration
	var limit int
	var interval, length string

	t := "time"
	if duration == time.Duration(0) {
		t = "limit"
		limit = qm.NumberPoints
		if qm.NumberPoints > 1000 {
			limit = 1000
		}
	} else {
		interval = fmt.Sprintf("%d seconds", int(duration.Seconds())/qm.NumberPoints)

		length = fmt.Sprintf("%d seconds", int(duration.Seconds()))
	}

	return &Query{
		Message:         qm,
		Type:            t,
		QueryUUID:       qm.QueryUUID,
		QueryField:      qm.QueryField,
		Limit:           limit,
		TimeInterval:    interval,
		Interval:        time.Duration(int(duration.Seconds()) / qm.NumberPoints * int(time.Second)),
		TimeLength:      length,
		ResponseChannel: make(chan *messages.QueryResponse, 1),
	}
}

/**
query historical data of uuid from startTime to endTime
prob should make sure they don't query like 1000 years or something
*/
func MakeQuery(qm *messages.QueryMessage) *Query {
	q := BuildQuery(qm)
	if qm.UseCache {
		queryCacheLock.Acquire("make-query")
		defer queryCacheLock.Release()
		hash := makeQueryHash(q)
		cacheItem, ok := queryCache[hash]
		if ok {
			if time.Since(cacheItem.lastUpdateTime) > qm.CacheDuration.Duration {
				q.ResponseChannel <- cacheItem.response
				cacheItem.lastUseTime = time.Now()
				return q
			}
		}

	}

	go makeQuery(q)
	return q
}

func makeQuery(query *Query) {
	val, exists := utils.GetVal(query.QueryUUID)
	if !exists {
		failedQuery(query, errors.New("uuid does not exist"))
	}

	var vals [][]interface{}
	var err error

	switch v := val.(type) {
	case *valuable.Stock:
		switch query.Type {
		case "time":
			vals, err = database.MakeStockHistoryTimeQuery(query.QueryUUID, query.TimeLength, query.QueryField, query.TimeInterval)
		case "limit":
			vals, err = database.MakeStockHistoryLimitQuery(query.QueryUUID, query.QueryField, query.Limit)
		}
	case *portfolio.Portfolio:
		switch query.Type {
		case "time":
			vals, err = database.MakePortfolioHistoryTimeQuery(query.QueryUUID, query.TimeLength, query.QueryField, query.TimeInterval)
		case "limit":
			vals, err = database.MakePortfolioHistoryLimitQuery(query.QueryUUID, query.QueryField, query.Limit)
		}
	case *ledger.Entry:
		switch query.Type {
		case "time":
			vals, err = database.MakeLedgerHistoryTimeQuery(query.QueryUUID, query.TimeLength, query.QueryField, query.TimeInterval)
		case "limit":
			vals, err = database.MakeLedgerHistoryLimitQuery(query.QueryUUID, query.QueryField, query.Limit)
		}

	default:
		fmt.Printf("%T", v)
	}

	if err != nil {
		failedQuery(query, err)
	}

	successQuery(query, vals)
}

func failedQuery(query *Query, err error) {
	query.ResponseChannel <- &messages.QueryResponse{
		Success: false,
		Error:   err.Error(),
		Message: query.Message,
	}
}

func successQuery(query *Query, values [][]interface{}) {
	queryCacheLock.Acquire("success-query")
	defer queryCacheLock.Release()

	response := &messages.QueryResponse{
		Success: true,
		Error:   "",
		Points:  values,
		Message: query.Message,
	}

	if query.Type == "time" {
		hash := makeQueryHash(query)
		_, ok := queryCache[hash]
		if !ok {
			queryCache[hash] = &queryCacheItem{
				lastUpdateTime: time.Now(),
				lastUseTime:    time.Now(),
				validTime:      query.Interval,
				response:       response,
				query:          query,
			}

		} else {
			queryCache[hash].lastUpdateTime = time.Now()
			queryCache[hash].response.Points = values

		}
	}
	query.ResponseChannel <- response

}
