package config

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Database {
    // Load environment variables
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Get MongoDB URI and database name from environment variables
    mongoURI := os.Getenv("MONGO_URI")
    dbName := os.Getenv("DB_NAME")

    // Connect to MongoDB
    client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal("Failed to create MongoDB client:", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Ping the database
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }

    log.Println("Connected to MongoDB")
    return client.Database(dbName)
}