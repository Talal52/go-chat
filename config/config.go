package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func ConnectDB() *sql.DB {
	connStr := "postgres://postgres:12345@localhost/chat_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Ping error:", err)
	}
	log.Println("Connected to DB")
	return db
}
