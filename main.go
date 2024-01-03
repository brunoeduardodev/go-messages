package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/brunoeduardodev/go-messages/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Json = map[string]interface{}
type ClientMessage struct {
	Message  string `json:"message"`
	Username string `json:"username"`
}

type Message struct {
	Content  string
	Username string
}

var messages = []Message{}
var connections = []*websocket.Conn{}

func getChatHTML() string {
	var messageBuffer bytes.Buffer
	template := utils.Must(template.ParseFiles("templates/chat-room.html"))
	template.Execute(&messageBuffer, Json{"Messages": messages})
	return messageBuffer.String()
}

func getSendMessageHTML() string {
	var messageBuffer bytes.Buffer
	template := utils.Must(template.ParseFiles("templates/send-message-form.html"))
	template.Execute(&messageBuffer, nil)
	return messageBuffer.String()

}

func main() {
	usernames := map[string]string{}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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

		connections = append(connections, conn)

		chatHtml := getChatHTML()
		conn.WriteMessage(1, []byte(chatHtml))

		for {
			_, rawMessage, err := conn.ReadMessage()
			if err != nil {
				return
			}

			var clientMessage ClientMessage
			err = json.Unmarshal(rawMessage, &clientMessage)
			if err != nil {
				fmt.Printf("Invalid message: %v\n", err)
				return
			}

			if clientMessage.Username != "" {
				usernames[conn.LocalAddr().String()] = clientMessage.Username

				sendMessageForm := getSendMessageHTML()
				conn.WriteMessage(1, []byte(sendMessageForm))
				continue
			}

			if usernames[conn.LocalAddr().String()] == "" {
				continue
			}

			messages = append(messages, Message{Content: clientMessage.Message, Username: usernames[conn.LocalAddr().String()]})
			for _, connection := range connections {
				chatHtml := getChatHTML()
				connection.WriteMessage(1, []byte(chatHtml))
			}
		}

	})

	fmt.Printf("🚀 Server started at http://localhost:8000/ \n")
	http.ListenAndServe(":8000", nil)
}
