package service

import (
	"errors"
	"log"
	"time"

	"github.com/Talal52/go-chat/chat/db"
	"github.com/Talal52/go-chat/chat/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *db.UserRepository
}

func NewAuthService(repo *db.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Signup(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return err
	}

	user := db.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.Repo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repo.GetUserByUsername(email) // Adjust to GetUserByEmail if renamed
	if err != nil {
		log.Println("User not found:", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Invalid password:", err)
		return "", err
	}

	jwtSecret := "mysecretkey"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Println("Error generating token:", err)
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) AuthenticateUser(email, password string) (*models.User, error) {
	user, err := s.Repo.GetUserByEmail(email)
	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Invalid password:", err)
		return nil, errors.New("invalid password")
	}

	modelsUser := &models.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}

	return modelsUser, nil
}
