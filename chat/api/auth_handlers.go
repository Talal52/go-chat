package api

import (
    "encoding/json"
    "github.com/Talal52/go-chat/chat/models"
    "github.com/Talal52/go-chat/chat/service"
    "net/http"
)

type AuthHandler struct {
    Service *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
    return &AuthHandler{Service: svc}
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
    var user models.User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if err := h.Service.Signup(user); err != nil {
        http.Error(w, "Could not create user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    token, err := h.Service.Login(credentials.Username, credentials.Password)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}