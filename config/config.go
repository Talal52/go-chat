package config

import (
    "context"
    "database/sql"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
    if err := godotenv.Load(".env"); err != nil {
        log.Fatal("Error loading .env file")
    }
}

type Config struct {
    HTTPPort     string
    WebSocketURL string
    PostgresURI  string
    MongoURI     string
    DBName       string
    JWTSecret    string
}

func LoadConfig() *Config {
    return &Config{
        HTTPPort:     os.Getenv("HTTP_PORT"),
        WebSocketURL: os.Getenv("WS_URL"),
        PostgresURI:  os.Getenv("POSTGRES_URI"),
        MongoURI:     os.Getenv("MONGO_URI"),
        DBName:       os.Getenv("DB_NAME"),
        JWTSecret:    os.Getenv("JWT_SECRET"),
    }
}

func ConnectPostgres(uri string) *sql.DB {
    db, err := sql.Open("postgres", uri)
    if err != nil {
        log.Fatalf("Failed to connect to PostgreSQL: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping PostgreSQL: %v", err)
    }

    log.Println("Connected to PostgreSQL")
    return db
}

func ConnectDB(uri, dbName string) *mongo.Database {
    clientOpts := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(context.Background(), clientOpts)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }

    log.Println("Connected to MongoDB")
    return client.Database(dbName)
}