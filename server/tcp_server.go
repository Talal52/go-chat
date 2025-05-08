package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Talal52/go-chat/chat/models"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/golang-jwt/jwt"
)

var (
	clients      = make(map[net.Conn]string)
	clientsMutex sync.Mutex
)

func StartTCPServer(chatService *service.ChatService) {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("TCP Server error: %v", err)
	}
	defer listener.Close()
	log.Println("TCP server started on :9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go handleConnection(conn, chatService)
	}
}

func handleConnection(conn net.Conn, chatService *service.ChatService) {
	defer conn.Close()

	if err := promptForToken(conn); err != nil {
		log.Println(err)
		return
	}

	name, err := authenticate(conn)
	if err != nil {
		log.Println(err)
		return
	}

	registerClient(conn, name)
	defer unregisterClient(conn, name)

	broadcast(fmt.Sprintf("%s joined the chat\n", name), conn)

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client disconnected: %v", conn.RemoteAddr())
			break
		}
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		if err := chatService.SaveMessage(models.Message{
			Sender:    name,
			Content:   msg,
			CreatedAt: time.Now(),
		}); err != nil {
			log.Printf("Error saving message: %v", err)
		}

		broadcast(fmt.Sprintf("%s: %s\n", name, msg), conn)
	}

	broadcast(fmt.Sprintf("%s left the chat\n", name), conn)
}

func promptForToken(conn net.Conn) error {
	_, err := conn.Write([]byte("Enter your JWT token: "))
	return err
}

func authenticate(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	tokenString, _ := reader.ReadString('\n')
	tokenString = strings.TrimSpace(tokenString)

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	name, ok := (*claims)["username"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token claims: username not found")
	}

	return name, nil
}

func registerClient(conn net.Conn, name string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	clients[conn] = name
}

func unregisterClient(conn net.Conn, name string) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	delete(clients, conn)
}

func broadcast(message string, sender net.Conn) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message))
		}
	}
}
