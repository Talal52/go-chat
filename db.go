package main

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("postgres", "postgres://postgres:12345@localhost/chat_db?sslmode=disable")
    if err != nil {
        log.Fatal("Database connection error:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Database ping error:", err)
    }
    log.Println("Connected to PostgreSQL database.")
}

func SaveMessage(sender, content string) {
    _, err := DB.Exec(`INSERT INTO messages (sender, content) VALUES ($1, $2)`, sender, content)
    if err != nil {
        log.Println("Error saving message:", err)
    }
}
