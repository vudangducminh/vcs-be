package main

import (
	"log"
	"net/http"

	// Add other imports here as needed

	// !! IMPORTANT: You will need to uncomment and update this line after running swag init !!
	// This imports the generated swagger documentation code.
	// Replace 'your_module_path' with the path from your go.mod file.
	// For example, if your go.mod is 'module github.com/yourname/sms', this would be:
	// _ "github.com/yourname/sms/docs"
	_ "sms/docs"

	httpSwagger "github.com/swaggo/http-swagger" // Import the http-swagger library
)

// @title SMS Project API
// @version 1.0
// @description API Documentation for the SMS project.
// @host localhost:8080 // <--- Change this to your actual host and port
// @BasePath /api/v1 // <--- Change this to your API's base path if you have one (e.g., /api/v1)
// @schemes http https // Specify the schemes your API supports
// @produce json // Default content type for responses
// @consumes json // Default content type for requests
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization // Example: For JWT or API key in Authorization header
// @securityDefinitions.basic BasicAuth // Example: For Basic Authentication

func main() {
	// ... your existing code for database connections, routing, etc. ...

	// Example of setting up a simple handler if you don't have a router yet
	// http.HandleFunc("/some_endpoint", yourHandlerFunction)

	// !! Add the handler for serving Swagger UI !!
	// This line mounts the Swagger UI at the /swagger/ path
	http.Handle("/swagger/", httpSwagger.Handler(
		// Provide the URL to the generated doc.json file
		// This URL should be accessible based on where you mount the http-swagger handler
		// If your API is served at localhost:8080 and swagger handler at /swagger/,
		// the doc.json is usually available at localhost:8080/swagger/doc.json
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // <--- Update host/port if needed
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"), // or "list" or "full"
		httpSwagger.DomID("swagger-ui"),
	))

	log.Println("Server starting on :8080...")                                      // <--- Update this to your actual port
	log.Println("Swagger UI available at http://localhost:8080/swagger/index.html") // <--- Update this

	// Start your HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil)) // <--- Update this to your actual port
}
