package db

import (
    "database/sql"
    "errors"
    "github.com/Talal52/go-chat/chat/models"
)

type AuthRepository struct {
    DB *sql.DB
}
type UserRepository struct {
    DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user models.User) error {
    _, err := r.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
    return err
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
    var user models.User
    err := r.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
    if err == sql.ErrNoRows {
        return nil, errors.New("user not found")
    }
    return &user, err
}