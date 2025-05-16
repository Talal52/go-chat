package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo; tighten in prod
	},
}

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

var clients = make(map[string]*Client) // Map userID to client

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	// Read userID from query (for demo)
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		log.Println("Missing userId")
		conn.Close()
		return
	}

	client := &Client{
		ID:   userID,
		Conn: conn,
		Send: make(chan []byte),
	}
	clients[userID] = client

	go writePump(client)
	readPump(client)
}

func readPump(client *Client) {
	defer func() {
		client.Conn.Close()
		delete(clients, client.ID)
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		// Here you decode message JSON and send to receiver client if connected
		// For demo, assume message is JSON with receiverID field (simple)
		// This logic can be extended

		receiverID := parseReceiverID(message)
		if receiverClient, ok := clients[receiverID]; ok {
			receiverClient.Send <- message
		} else {
			log.Println("Receiver not connected:", receiverID)
		}
	}
}

func writePump(client *Client) {
	for {
		msg, ok := <-client.Send
		if !ok {
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		err := client.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func parseReceiverID(message []byte) string {
	// You can decode message JSON and return receiverID
	// For demo, hardcode or simple parse (needs improvement)
	return "receiverUserID"
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	log.Println("WebSocket Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
