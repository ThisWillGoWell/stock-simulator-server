package web

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/config"

	metrics "github.com/ThisWillGoWell/stock-simulator-server/src/metics"

	"github.com/ThisWillGoWell/stock-simulator-server/src/account"
	"github.com/ThisWillGoWell/stock-simulator-server/src/client"
	"github.com/ThisWillGoWell/stock-simulator-server/src/messages"
	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]http.Client) // connected clients

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func StartHandlers() {
	//shareDir := os.Getenv("FILE_SERVE")
	//if shareDir == "" {
	//	shareDir = "static"
	//}
	//fmt.Println(shareDir)
	//var fs = http.FileServer(http.Dir(shareDir))

	http.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
		config.Seed()
		http.Redirect(w, r, "/", 301)
	})

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
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
		token, err := account.NewUser(pair[0], displayName, pair[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		io.WriteString(w, token)
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
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
		token, err := account.ValidateUser(pair[0], pair[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		io.WriteString(w, token)
	})

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		b, _ := json.Marshal(metrics.Counter)
		io.WriteString(w, string(b))
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func ServePath(p string) {
	var fs = http.FileServer(http.Dir(p))
	http.Handle("/", fs)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "DisplayName, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got Upgrade")
	//first upgrade the connection
	ws, err := upgrader.Upgrade(w, r, nil)
	defer ws.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	socketRX := make(chan string)
	socketTX := make(chan string)
	// Gate Keeper
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}
		loginErr := client.InitialReceive(string(msg), socketTX, socketRX)
		if loginErr != nil {
			val, err := json.Marshal(messages.FailedConnect(loginErr))
			if err != nil {
				fmt.Print(err)
			}
			ws.WriteMessage(websocket.TextMessage, val)
			// Give time for the connection to send the error
			<-time.After(100 * time.Millisecond)
			return
		} else {
			break
		}

	}
	// Make sure we close the connection when the function returns

	go runTxSocket(ws, socketTX)
	rxSocket(ws, socketRX)
}

func runTxSocket(conn *websocket.Conn, tx chan string) {
	for str := range tx {
		conn.WriteMessage(websocket.TextMessage, []byte(str))
	}
}

func rxSocket(conn *websocket.Conn, rx chan string) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		rx <- string(msg)
	}
	close(rx)
}
