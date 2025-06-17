package websocket

import (
	"log"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	conn *websocket.Conn
	url  string
}

func NewClient(url string) *WSClient {
	return &WSClient{url: url}
}

func (c *WSClient) Connect() error {
	var dialer websocket.Dialer
	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *WSClient) ListenMessages() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		log.Printf("Received: %s", message)
	}
}
