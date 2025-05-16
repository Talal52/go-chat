package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/Talal52/go-chat/server"
    "github.com/Talal52/go-chat/config"
)

func main() {
    cfg := config.LoadConfig()
    srv := server.NewHTTPServer(cfg)
    fmt.Printf("HTTP Sender running on port %s\n", cfg.HTTPPort)
    log.Fatal(http.ListenAndServe(":"+cfg.HTTPPort, srv))
}
