package http

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/ThisWillGoWell/stock-simulator-server/src/objects/user"
)

func makeUser(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}

	if r.Method != "PUT" {
		http.Error(w, "put only", http.StatusMethodNotAllowed)
		return
	}

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	displayName := r.Header.Get("DisplayName")

	if len(auth) != 2 || auth[0] != "Basic" {
		http.Error(w, "create failed", http.StatusBadRequest)
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		http.Error(w, "create failed", http.StatusBadRequest)
		return
	}
	token, err := user.NewUser(pair[0], displayName, pair[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	io.WriteString(w, token)
}

func getToken(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		return
	}
	if r.Method != "GET" {
		http.Error(w, "get only", http.StatusMethodNotAllowed)
		return
	}
	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		http.Error(w, "authorization failed", http.StatusUnauthorized)
		return
	}
	token, err := user.ValidateUser(pair[0], pair[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	io.WriteString(w, token)
}
