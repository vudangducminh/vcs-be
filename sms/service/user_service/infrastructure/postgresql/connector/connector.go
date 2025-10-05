package connector

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var PostgreSQL *xorm.Engine
var isConnected = false

func IsConnected() bool {
	return isConnected
}

func Connect() {
	// Get database configuration from environment variables with fallbacks
	host := getEnv("POSTGRES_HOST", "postgresql")
	port := getEnv("POSTGRES_PORT", "5432")
	user := getEnv("POSTGRES_USER", "vudangducminh")
	password := getEnv("POSTGRES_PASSWORD", "Amogus69420")
	dbname := getEnv("POSTGRES_DB", "postgres")

	conns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	log.Printf("Attempting to connect to PostgreSQL at %s:%s", host, port)

	var err error
	PostgreSQL, err = xorm.NewEngine("postgres", conns)
	if err != nil {
		log.Printf("Failed to connect to PostgreSQL: %v", err)
		log.Fatal(err)
		return
	}

	// Test the connection
	if err = PostgreSQL.Ping(); err != nil {
		log.Printf("Failed to ping PostgreSQL: %v", err)
		log.Fatal(err)
		return
	}

	isConnected = true
	log.Println("Connected to PostgreSQL database successfully")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
