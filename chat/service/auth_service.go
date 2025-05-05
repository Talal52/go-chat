package service

import (
    "errors"
    "github.com/Talal52/go-chat/chat/db"
    "github.com/Talal52/go-chat/chat/models"
    "golang.org/x/crypto/bcrypt"
    "time"

    "github.com/golang-jwt/jwt"
)

var secretKey = []byte("secret-key")

type AuthService struct {
    Repo *db.AuthRepository
}

func NewAuthService(repo *db.AuthRepository) *AuthService {
    return &AuthService{Repo: repo}
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
    return token.SignedString(secretKey)
}