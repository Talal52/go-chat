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

func ConnectPostgres() *sql.DB {
	connStr := os.Getenv("POSTGRES_URI")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	return db
}

func ConnectDB() *mongo.Database {
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	clientOpts := options.Client().ApplyURI(mongoURI)
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
type Config struct {
	HTTPPort     string
	WebSocketURL string
	JWTSecret    string
}

func LoadConfig() *Config {
	return &Config{
		HTTPPort:     getEnv("HTTP_PORT", "8082"),
		WebSocketURL: getEnv("WS_URL", "ws://localhost:8081/ws"),
		JWTSecret:    getEnv("JWT_SECRET", "mysecretkey"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
