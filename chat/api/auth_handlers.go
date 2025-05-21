package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/Talal52/go-chat/chat/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthHandler struct {
	Service *service.AuthService
}

func (h *AuthHandler) SignupGin(c *gin.Context) {
	log.Println("Received signup request")
	var user struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Invalid input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.Service.Signup(user.Email, user.Password); err != nil {
		log.Println("Signup error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (h *AuthHandler) LoginGin(c *gin.Context) {
	log.Println("Received login request")
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Token    string `json:"token"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Println("Invalid input:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if credentials.Token != "" {
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "mysecretkey"
		}

		token, err := jwt.Parse(credentials.Token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid or expired JWT token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired JWT token"})
			return
		}
	}

	newToken, err := h.Service.Login(credentials.Email, credentials.Password)
	if err != nil {
		log.Println("Login error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
