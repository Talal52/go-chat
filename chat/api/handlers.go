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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHandler struct {
	Service *service.ChatService
}

func NewChatHandler(svc *service.ChatService) *ChatHandler {
	return &ChatHandler{Service: svc}
}

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

func (h *ChatHandler) GetMessagesGin(c *gin.Context) {
	messages, err := h.Service.GetMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving messages"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

// HTTP endpoint to fetch group messages
func (h *ChatHandler) GetGroupMessages(c *gin.Context) {
	groupIDStr := c.Query("group_id")
	groupID, err := primitive.ObjectIDFromHex(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid group ID"})
		return
	}

	messages, err := h.Service.GetMessagesByGroupID(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

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

