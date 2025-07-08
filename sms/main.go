package main

import (
	"log"
	"sms/API/swagger"
	_ "sms/API/users_handler" // Importing users_handler for Swagger documentation
	_ "sms/docs"
	"sms/server/database/postgresql/connector"
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

func main() {
	// Initialize the database connection
	connector.ConnectToDB()
	if !connector.IsConnected() {
		log.Println("Failed to connect to the database")
	}

	// Connect to Swagger for API documentation
	swagger.ConnectToSwagger()

	log.Println("Server starting on http://localhost:8800")
	log.Println("Swagger UI available at http://localhost:8800/swagger/index.html")

}
