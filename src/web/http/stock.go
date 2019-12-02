package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/valuable"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects"
	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/user"
)

type NewStockRequest struct {
	Stock       objects.Stock `json:"object"`
	AccessToken string        `json:"access_token"`
}

func makeStock(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	if (*r).Method == "OPTIONS" {
		return
	}

	if r.Method != "POST" {
		http.Error(w, "put only", http.StatusMethodNotAllowed)
		return
	}

	newStock, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request := NewStockRequest{}
	if err := json.Unmarshal(newStock, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := user.GetUserFromToken(request.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if !u.IsAdmin {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = valuable.NewStock(request.Stock.TickerId, request.Stock.Name, request.Stock.CurrentPrice, 10*time.Second)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
