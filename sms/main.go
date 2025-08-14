package main

import (
	"log"
	"sms/API/swagger"
	_ "sms/API/users_handler" // Importing users_handler for Swagger documentation
	_ "sms/docs"
	"sms/object"
	redis "sms/server/database/cache/redis/connector"
	elastic "sms/server/database/elasticsearch/connector"
	elastic_query "sms/server/database/elasticsearch/query"
	postgresql "sms/server/database/postgresql/connector"
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
	// go time_checker.TimeCheckerForSendingEmails()

	// Connect to Swagger for API documentation
	swagger.ConnectToSwagger()
	log.Println("Server starting on http://localhost:8800")
	log.Println("Swagger UI available at http://localhost:8800/swagger/index.html")
}

func GenerateServer() {
	var servers []object.Server
	var sv1 = object.Server{
		ServerId:        "11244475057891453f598530bd8f2b702dd56bdfa125d10f30c288ff5529ed47",
		ServerName:      "Server 9",
		Status:          "inactive",
		Uptime:          3600, // 1 hour in seconds
		CreatedTime:     1755140989,
		LastUpdatedTime: 1755140989,
		IPv4:            "192.168.1.9",
	}
	servers = append(servers, sv1)
	// for i := 1; i < 10000; i++ {
	// 	var rng = rand.Intn(1000) % 3 // Random number between 1 and 1000
	// 	var status string
	// 	switch rng {
	// 	case 0:
	// 		status = "active"
	// 	case 1:
	// 		status = "inactive"
	// 	case 2:
	// 		status = "maintenance"
	// 	default:
	// 		status = "other"
	// 	}
	// 	server := object.Server{
	// 		ServerId:        algorithm.SHA256Hash(time.Now().String() + fmt.Sprintf("%d", i)),
	// 		ServerName:      "Server " + fmt.Sprintf("%d", i),
	// 		Status:          status,
	// 		IPv4:            "192.168.1." + fmt.Sprintf("%d", i),
	// 		Uptime:          3600, // 1 hour in seconds
	// 		CreatedTime:     time.Now().Unix(),
	// 		LastUpdatedTime: time.Now().Unix(),
	// 	}
	// 	servers = append(servers, server)
	// }
	status := elastic_query.BulkServerInfo(servers)
	if status != 201 {
		log.Println("Failed to add servers to Elasticsearch in bulk, status code:", status)
		return
	}
	log.Println("Generated and added 10000 sample servers to Elasticsearch")
}
