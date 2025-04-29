package main

import (
	"github.com/Talal52/go-chat/config"
	"github.com/Talal52/go-chat/server"
)

func main() {
	dbConn := config.ConnectDB()
	server.InitServers(dbConn)

	select {} // block forever
}
