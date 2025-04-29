package db

import (
	"database/sql"
	"github.com/Talal52/go-chat/chat/models"
)

type ChatRepository struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{DB: db}
}

func (r *ChatRepository) SaveMessage(msg models.Message) error {
	_, err := r.DB.Exec(`INSERT INTO messages (sender, content, created_at) VALUES ($1, $2, $3)`,
		msg.Sender, msg.Content, msg.CreatedAt)
	return err
}

func (r *ChatRepository) GetMessages() ([]models.Message, error) {
	rows, err := r.DB.Query(`SELECT sender, content, created_at FROM messages ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.Sender, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
