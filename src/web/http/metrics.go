package http

import (
	"encoding/json"
	"io"
	"net/http"

	metrics "github.com/ThisWillGoWell/stock-simulator-server/src/app/metics"
)

func getMetrics(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	b, _ := json.Marshal(metrics.Counter)
	io.WriteString(w, string(b))
}
