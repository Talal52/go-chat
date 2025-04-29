package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:9000")
    if err != nil {
        fmt.Println("Connection error:", err)
        return
    }
    defer conn.Close()

    go receiveMessages(conn)

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text()
        if text != "" {
            fmt.Fprintf(conn, text+"\n")
        }
    }
}

func receiveMessages(conn net.Conn) {
    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("Disconnected from server.")
            return
        }
        fmt.Print(msg)
    }
}
