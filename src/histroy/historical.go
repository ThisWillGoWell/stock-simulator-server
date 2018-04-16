package histroy

import (
	"errors"
	"github.com/stock-simulator-server/src/portfolio"
	"github.com/stock-simulator-server/src/utils"
	"github.com/stock-simulator-server/src/valuable"
	"time"
)

/**
query historical data of uuid from startTime to endTime
prob should make sure they don't query like 1000 years or something
*/
func Query(uuid string, startTime, endTime time.Time) (*TimeSeriesObject, error) {
	val, exists := utils.GetVal(uuid)
	if !exists {
		return nil, errors.New("uuid does not exist")
	}
	switch val.(type) {
	case valuable.Stock:
		return queryStock(uuid, startTime, endTime), nil
	case portfolio.Portfolio:
		return queryPortfolio(uuid, startTime, endTime), nil
	}
	return nil, errors.New("something went wrong, query switch failed")
}

func queryPortfolio(uuid string, startTime, endtime time.Time) *TimeSeriesObject {
	testTs := make([]*TimeSeriesEntry, 10)
	var i int64
	for i = 0; i < 10; i++ {
		testTs[i] = &TimeSeriesEntry{
			Timestamp: time.Unix(i, 0),
			Value:     i,
		}
	}

	return &TimeSeriesObject{
		uuid,
		testTs,
	}
}

func queryStock(uuid string, startTime, endtime time.Time) *TimeSeriesObject {
	return nil
}

type TimeSeriesObject struct {
	Uuid string `json:"uuid"`
	Data []*TimeSeriesEntry
}

type TimeSeriesEntry struct {
	Timestamp time.Time
	Value     interface{}
}
