package main

// import (
//     "encoding/json"
//     "net/http"
//     "strconv"
// )

// type Message struct {
//     Sender    string `json:"sender"`
//     Content   string `json:"content"`
//     CreatedAt string `json:"created_at"`
// }

// func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
//     pageStr := r.URL.Query().Get("page")
//     pageSizeStr := r.URL.Query().Get("pageSize")

//     page, err := strconv.Atoi(pageStr)
//     if err != nil || page <= 0 {
//         page = 1
//     }

//     pageSize, err := strconv.Atoi(pageSizeStr)
//     if err != nil || pageSize <= 0 {
//         pageSize = 10
//     }

//     messages, err := FetchMessages(page, pageSize)
//     if err != nil {
//         http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(messages)
// }

// func FetchMessages(page, pageSize int) ([]Message, error) {
//     offset := (page - 1) * pageSize

//     rows, err := DB.Query(`SELECT sender, content, created_at FROM messages ORDER BY created_at DESC LIMIT $1 OFFSET $2`, pageSize, offset)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()

//     var messages []Message

//     for rows.Next() {
//         var msg Message
//         if err := rows.Scan(&msg.Sender, &msg.Content, &msg.CreatedAt); err != nil {
//             return nil, err
//         }
//         messages = append(messages, msg)
//     }

//     return messages, nil
// }
