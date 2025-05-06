package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Talal52/go-chat/chat/models"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/golang-jwt/jwt"
)

type ChatHandler struct {
	Service *service.ChatService
}

func NewChatHandler(svc *service.ChatService) *ChatHandler {
	return &ChatHandler{Service: svc}
}

func (h *ChatHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	// Extract JWT token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("Authorization header missing or invalid")
		http.Error(w, "Unauthorized: Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")


	// Parse and validate the token
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		secret := os.Getenv("JWT_SECRET")
		log.Println("Using JWT_SECRET:", secret)
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		log.Printf("Invalid token: %v", err) // Log the error
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract the sender's username from the token claims
	sender, ok := (*claims)["username"].(string)
	if !ok {
		log.Println("Invalid token claims: username not found") // Log the issue
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	log.Println("Sender extracted from token:", sender)

	// Decode the message from the request body
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		log.Println("Error decoding message body:", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set the sender and timestamp
	msg.Sender = sender
	msg.CreatedAt = time.Now()

	// Save the message using the service
	if err := h.Service.SaveMessage(msg); err != nil {
		log.Printf("Error saving message: %v", err) // Log the error
		http.Error(w, "Could not save message", http.StatusInternalServerError)
		return
	}

	log.Println("Message saved successfully")
	w.WriteHeader(http.StatusCreated)
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	messages, err := h.Service.GetMessages()
	if err != nil {
		http.Error(w, "Error retrieving messages", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}
