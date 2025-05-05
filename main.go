package main

import (
    "github.com/Talal52/go-chat/config"
    "github.com/Talal52/go-chat/server"
)

func main() {
    // Connect to MongoDB
    mongoDB := config.ConnectDB()

    // Connect to PostgreSQL
    postgresDB := config.ConnectPostgres()

    // Initialize servers
    server.InitServers(mongoDB, postgresDB)

    // Block forever
    select {}
}