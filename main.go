package main

import (
    "log"
    "net/http"
    "sync"
)

func main() {
    // Initialize Database
    InitDB()
    defer DB.Close()

    var wg sync.WaitGroup

    wg.Add(2)

    // Start TCP Server (chat server)
    go func() {
        defer wg.Done()
        StartTCPServer()
    }()

    // Start HTTP Server (pagination API)
    go func() {
        defer wg.Done()
        http.HandleFunc("/messages", GetMessagesHandler)
        log.Println("HTTP server started on :8000")
        log.Fatal(http.ListenAndServe(":8000", nil))
    }()

    wg.Wait()
}
