package histroy

func QueryHistroy(string uuid) {

}

type TimeSeriesObject struct {
	Uuid string `json:"uuid"`
	data []*TimeSeriesObject
}

type TimeSeriesEntry struct {
	Timestamp int64
	Value     interface{}
}
