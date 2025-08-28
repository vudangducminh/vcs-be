package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sms/algorithm"
	_ "sms/docs"
	"sms/object"
	redis "sms/server/database/cache/redis/connector"
	elastic "sms/server/database/elasticsearch/connector"
	elastic_query "sms/server/database/elasticsearch/query"
	postgresql "sms/server/database/postgresql/connector"
	time_checker "sms/service/report_service/time_checker"
	"sms/service/swagger"
	_ "sms/service/user_service" // Importing users_handler for Swagger documentation
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
// @Tag.name		Servers
// @Tag.description "Operations related to server management"
func main() {
	// Set the log file
	file, err := os.OpenFile("log/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	// Set the output of the default logger to the file.
	// All subsequent calls to log.Println, log.Printf, etc. will write to this file.
	log.SetOutput(file)
	log.Println("Logger initialized. Subsequent logs will be written to log/server.log")

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
	// GenerateServer()

	// Start the time checker for daily report email requests
	go time_checker.TimeCheckerForSendingEmails()

	// Connect to Swagger for API documentation
	swagger.ConnectToSwagger()
	log.Println("Server starting on http://localhost:8800")
	log.Println("Swagger UI available at http://localhost:8800/swagger/index.html")
}

func GenerateServer() {
	var servers []object.Server
	for i := 0; i < 10000; i++ {
		var rng = rand.Intn(1000)
		if rng < 600 {
			rng = 0 // Active status
		} else if rng < 900 {
			rng = 1 // Inactive status
		} else {
			rng = 2 // Maintenance status
		}
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
			IPv4:            fmt.Sprintf("%d", rand.Intn(256)) + "." + fmt.Sprintf("%d", rand.Intn(256)) + "." + fmt.Sprintf("%d", rand.Intn(256)) + "." + fmt.Sprintf("%d", rand.Intn(256)),
			Uptime:          int64(rand.Intn(86400)),                               // 1 hour in seconds
			CreatedTime:     time.Now().Unix() - 86400*2 - int64(rand.Intn(86400)), // Created 2 to 3 days ago
			LastUpdatedTime: time.Now().Unix(),
		}
		servers = append(servers, server)
	}
	status := elastic_query.BulkServerInfo(servers)
	if status != 201 {
		log.Println("Failed to add servers to Elasticsearch in bulk, status code:", status)
		return
	}
	log.Println("Generated and added 10000 sample servers to Elasticsearch")
}
