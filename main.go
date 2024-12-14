package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func reader(conn *websocket.Conn) {
	defer conn.Close() // Ulanishni yopishni unutmang
	for {
		// Xabarni o'qish
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		// Olingan xabarni chop etish
		log.Println("Received:", string(p))

		// Xabarni qayta yozish (send back)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// WebSocketga ulanishni yangilash
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	log.Println("Client Connected")

	// Xush kelibsiz xabarini yuborish
	err = ws.WriteMessage(websocket.TextMessage, []byte("Hi Client!"))
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}

	// Mijoz bilan aloqa qilish
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Server started on :8080")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
