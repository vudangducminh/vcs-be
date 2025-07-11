package main

import (
	"log"
	"sms/API/swagger"
	_ "sms/API/users_handler" // Importing users_handler for Swagger documentation
	_ "sms/docs"
	redis "sms/server/database/cache/redis/connector"
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
	// Connect to Swagger for API documentation
	swagger.ConnectToSwagger()
	log.Println("Server starting on http://localhost:8800")
	log.Println("Swagger UI available at http://localhost:8800/swagger/index.html")

}
