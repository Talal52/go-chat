package service

import (
    "errors"
    "time"

    "github.com/Talal52/go-chat/chat/db"
    "github.com/Talal52/go-chat/chat/models"
    "github.com/golang-jwt/jwt"
    "golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")

type AuthService struct {
    Repo *db.UserRepository
}

func NewAuthService(repo *db.UserRepository) *AuthService {
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