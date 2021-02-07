package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ThisWillGoWell/stock-simulator-server/src/messages"
	"github.com/ThisWillGoWell/stock-simulator-server/src/web/client"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnections(w http.ResponseWriter, r *http.Request) {

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
