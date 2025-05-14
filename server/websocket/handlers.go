package websocket

import (
    "net/http"
)

func WebSocketHandler(server *WebSocketServer) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        server.HandleConnections(w, r)
    }
}