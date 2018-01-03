package web

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)

var fs = http.FileServer(http.Dir("static"))

func StartHandlers() {
	http.Handle("/", fs)

	http.ListenAndServe(":3000", nil)

}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func startWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(conn)
		return
	}

}
