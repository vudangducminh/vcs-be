package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sms/API/swagger"
	_ "sms/API/users_handler" // Importing users_handler for Swagger documentation
	"sms/algorithm"
	_ "sms/docs"
	time_checker "sms/notification/time_checker"
	"sms/object"
	redis "sms/server/database/cache/redis/connector"
	elastic "sms/server/database/elasticsearch/connector"
	elastic_query "sms/server/database/elasticsearch/query"
	postgresql "sms/server/database/postgresql/connector"
	"time"
)

// @title           VCS System Management API
// @version         1.0
// @description     This is a sample server for VCS System Management API.
// @contact.name    Vu Dang Duc Minh
// @contact.email   vudangducminh@gmail.com
// @contact.url     https://github.com/vudangducminh
// @BasePath        /api/v1
// @schemes         http
// @host            localhost:8800
// @Tag.name        Users
// @Tag.description "Operations related to user authentication and management"
// @Tag.name		Auth
// @Tag.description "Operations related to user authentication"
// @Tag.name		Servers
// @Tag.description "Operations related to server management"
func main() {
	// Initialize the database connection
	postgresql.ConnectToDB()
	if !postgresql.IsConnected() {
		log.Println("Failed to connect to postgresql database")
	}
	// Initialize the Redis connection
	redis.ConnectToRedis()

	if !redis.IsConnected() {
		log.Println("Failed to connect to redis database")
	}
	// Initialize the Elasticsearch connection
	elastic.ConnectToEs()
	if !elastic.IsConnected() {
		log.Println("Failed to connect to Elasticsearch")
	}

	// Generate sample servers for testing
	GenerateServer()

	// Start the time checker for daily report email requests
	go time_checker.TimeCheckerForSendingEmails()
	go time_checker.CheckServerUptime()

	// Connect to Swagger for API documentation
	swagger.ConnectToSwagger()
	log.Println("Server starting on http://localhost:8800")
	log.Println("Swagger UI available at http://localhost:8800/swagger/index.html")
}

func GenerateServer() {
	for i := 0; i < 10000; i++ {
		var rng = rand.Intn(1000) % 3 // Random number between 1 and 1000
		var status string
		switch rng {
		case 0:
			status = "active"
		case 1:
			status = "inactive"
		case 2:
			status = "maintenance"
		default:
			status = "other"
		}
		server := object.Server{
			ServerId:        algorithm.SHA256Hash(time.Now().String() + fmt.Sprintf("%d", i)),
			ServerName:      "Server " + fmt.Sprintf("%d", i),
			Status:          status,
			IPv4:            "192.168.1." + fmt.Sprintf("%d", i),
			Uptime:          3600, // 1 hour in seconds
			CreatedTime:     time.Now().Format(time.RFC3339),
			LastUpdatedTime: time.Now().Format(time.RFC3339),
		}
		statusCode := elastic_query.AddServerInfo(server)
		if statusCode != http.StatusCreated {
			log.Printf("Failed to add server %s: %v", server.ServerName, status)
			continue
		}
		log.Println("Server added successfully:", server.ServerName)
	}
}
