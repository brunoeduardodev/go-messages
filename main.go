package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Hello</h1>"))
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("Error upgrading connection %v\n", err)
			w.WriteHeader(400)
			w.Write([]byte("{\"error:\" \"Something went wrong\"}"))
			return
		}

		conn.WriteMessage(1, []byte("Hello from server!"))
	})

	http.ListenAndServe(":8000", nil)
}
