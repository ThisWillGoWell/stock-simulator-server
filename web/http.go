package web

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)



var clients = make(map[*websocket.Conn]http.Client) // connected clients

func StartHandlers() {
	var fs = http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections
	http.ListenAndServe(":3000", nil)

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	username := r.Header["username"]
	password := r.Header["password"]
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()
}



func startWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(conn)
		return
	}

}
