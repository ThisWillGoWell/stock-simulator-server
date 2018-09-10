package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stock-simulator-server/src/account"
	"github.com/stock-simulator-server/src/app"
	"github.com/stock-simulator-server/src/client"
	"github.com/stock-simulator-server/src/messages"
	"log"
	"net/http"
	"time"
)

var clients = make(map[*websocket.Conn]http.Client) // connected clients

func StartHandlers() {
	//shareDir := os.Getenv("FILE_SERVE")
	//if shareDir == "" {
	//	shareDir = "static"
	//}
	//fmt.Println(shareDir)
	//var fs = http.FileServer(http.Dir(shareDir))

	//http.Handle("/", fs)
	http.HandleFunc("/load", func(w http.ResponseWriter, r *http.Request) {
		go app.LoadVars()
		<- time.After(time.Second)
		http.Redirect(w, r, "/", 200)
	})
	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) { account.NewUser("Will", "pass") })

	http.HandleFunc("/ws", handleConnections)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},

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
			val, err := json.Marshal(messages.FailedLogin(loginErr))
			if err != nil{
				fmt.Print(err)
			}
			ws.WriteMessage(websocket.TextMessage, val)
			// Give time for the connection to send the error
			<- time.After(100 * time.Millisecond)
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
