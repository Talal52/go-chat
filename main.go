package main

import (
    "github.com/Talal52/go-chat/config"
    "github.com/Talal52/go-chat/server"
)

func main() {
    // Connect to MongoDB
    dbConn := config.ConnectDB()

    // Initialize servers
    server.InitServers(dbConn)

    // Block forever
    select {}
}