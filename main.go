package main

import (
	"github.com/Talal52/go-chat/config"
	"github.com/Talal52/go-chat/server"
)

func main() {
	mongoDB := config.ConnectDB()

	postgresDB := config.ConnectPostgres()

	// Initialize servers
	server.InitServers(mongoDB, postgresDB)

	select {}
}
