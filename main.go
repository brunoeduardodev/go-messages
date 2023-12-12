package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/brunoeduardodev/go-messages/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := utils.Must(template.ParseFiles("templates/index.html"))
		template.Execute(w, nil)
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
