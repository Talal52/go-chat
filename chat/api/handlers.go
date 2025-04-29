package api

import (
	"encoding/json"
	"github.com/Talal52/go-chat/chat/models"
	"github.com/Talal52/go-chat/chat/service"
	"net/http"
	"time"
)

type ChatHandler struct {
	Service *service.ChatService
}

func NewChatHandler(svc *service.ChatService) *ChatHandler {
	return &ChatHandler{Service: svc}
}

func (h *ChatHandler) PostMessage(w http.ResponseWriter, r *http.Request) {
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	msg.CreatedAt = time.Now()
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
