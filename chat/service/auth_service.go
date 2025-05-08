package service

import (
	"errors"
	"os"
	"time"

	"github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo     *db.UserRepository
	secretKey []byte
}

func NewAuthService(repo *db.UserRepository) *AuthService {
	secretKey := []byte(os.Getenv("JWT_SECRET")) // Load the secret key once during initialization
	return &AuthService{Repo: repo, secretKey: secretKey}
}

func (s *AuthService) Signup(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.Repo.CreateUser(user)
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(s.secretKey)
}
