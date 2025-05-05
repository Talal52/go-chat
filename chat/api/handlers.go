package api

import (
	"encoding/json"
	"net/http"
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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate the token
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil // Replace with your secret key
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract the sender's username from the token claims
	sender, ok := (*claims)["username"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Decode the message from the request body
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Set the sender and timestamp
	msg.Sender = sender
	msg.CreatedAt = time.Now()

	// Save the message using the service
	if err := h.Service.SaveMessage(msg); err != nil {
		http.Error(w, "Could not save message", http.StatusInternalServerError)
		return
	}

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
