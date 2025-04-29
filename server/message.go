// filepath: c:\Users\PMYLS\Desktop\office\go-chat\server\message.go
package server

// import (
//     "database/sql"
//     "log"
// )

// var db *sql.DB // Global database connection (initialize this elsewhere)

// // SaveMessage saves a chat message to the database
// func SaveMessage(username, message string) {
//     query := "INSERT INTO messages (username, message) VALUES (?, ?)"
//     _, err := db.Exec(query, username, message)
//     if err != nil {
//         log.Println("Error saving message:", err)
//     }
// }