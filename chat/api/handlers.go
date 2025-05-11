package api

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Talal52/go-chat/chat/models"
	"github.com/Talal52/go-chat/chat/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type ChatHandler struct {
	Service *service.ChatService
}

func NewChatHandler(svc *service.ChatService) *ChatHandler {
	return &ChatHandler{Service: svc}
}

// PostMessageGin handles posting a message using Gin and JWT auth.
func (h *ChatHandler) PostMessageGin(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenString, err := extractToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
		return
	}

	sender, err := parseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var msg models.Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	msg.Sender = sender
	msg.CreatedAt = time.Now()

	if err := h.Service.SaveMessage(msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message posted successfully"})
}

// func (h *ChatHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
// 	// Extract and validate JWT token from Authorization header
// 	tokenString, err := extractToken(r.Header.Get("Authorization"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	// Parse and validate the token
// 	sender, err := parseToken(tokenString)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusUnauthorized)
// 		return
// 	}

// 	// Decode the message from the request body
// 	var msg models.Message
// 	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	// Set the sender and timestamp
// 	msg.Sender = sender
// 	msg.CreatedAt = time.Now()

// 	// Save the message using the service
// 	if err := h.Service.SaveMessage(msg); err != nil {
// 		http.Error(w, "Could not save message", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }

// func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
// 	messages, err := h.Service.GetMessages()
// 	if err != nil {
// 		http.Error(w, "Error retrieving messages", http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(messages)
// }

// GetMessagesGin handles retrieving messages using the Gin framework.
func (h *ChatHandler) GetMessagesGin(c *gin.Context) {
	messages, err := h.Service.GetMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// extractToken extracts the token from the Authorization header
func extractToken(authHeader string) (string, error) {
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return "", http.ErrNoCookie
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// parseToken parses the JWT token and extracts the sender's username
func parseToken(tokenString string) (string, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", http.ErrNoCookie
	}

	sender, ok := (*claims)["username"].(string)
	if !ok {
		return "", http.ErrNoCookie
	}
	return sender, nil
}
