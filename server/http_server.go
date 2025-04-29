package server

import (
	"github.com/Talal52/go-chat/chat/api"
	"log"
	"net/http"
)

func StartHTTPServer(handler *api.ChatHandler) {
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetMessages(w, r)
		} else if r.Method == http.MethodPost {
			handler.PostMessage(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Println("HTTP server running on :8000")
	http.ListenAndServe(":8000", nil)
}
