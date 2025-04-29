package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
    "sync"
)

var clients = make(map[net.Conn]string) // connection -> username
var clientsMutex sync.Mutex

func StartTCPServer() {
    listener, err := net.Listen("tcp", ":9000")
    if err != nil {
        log.Fatal("TCP Server error:", err)
    }
    defer listener.Close()
    log.Println("TCP server started on :9000")

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Connection error:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    conn.Write([]byte("Enter your name: "))
    reader := bufio.NewReader(conn)
    name, _ := reader.ReadString('\n')
    name = strings.TrimSpace(name)

    clientsMutex.Lock()
    clients[conn] = name
    clientsMutex.Unlock()

    broadcast(fmt.Sprintf("%s joined the chat\n", name), conn)

    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            log.Println("Client disconnected:", conn.RemoteAddr())
            break
        }
        msg = strings.TrimSpace(msg)
        if msg == "" {
            continue
        }

        SaveMessage(name, msg) // Save in DB
        broadcast(fmt.Sprintf("%s: %s\n", name, msg), conn)
    }

    clientsMutex.Lock()
    delete(clients, conn)
    clientsMutex.Unlock()

    broadcast(fmt.Sprintf("%s left the chat\n", name), conn)
}

func broadcast(message string, sender net.Conn) {
    clientsMutex.Lock()
    defer clientsMutex.Unlock()
    for conn := range clients {
        if conn != sender {
            conn.Write([]byte(message))
        }
    }
}
