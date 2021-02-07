package http

import (
	"log"
	"net/http"

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

	//http.HandleFunc("/api/load", func(w http.ResponseWriter, r *http.Request) {
	//	http.Redirect(w, r, "/", 301)
	//})

	http.HandleFunc("/api/stock", makeStock)
	http.HandleFunc("/api/create", makeUser)
	http.HandleFunc("/api/token", getToken)
	http.HandleFunc("/api/ws", handleConnections)
	http.HandleFunc("/api/metrics", getMetrics)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		w.WriteHeader(200)
	})

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "DisplayName, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func MakeStock() {

}
