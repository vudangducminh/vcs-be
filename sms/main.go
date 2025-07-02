package main

import (
	"log"
	"net/http"
	"sms/client/handler"
	"sms/database/postgresql/connector"
)

func main() {
	// InWitialize the database connection
	connector.ConnectToDB()
	if !connector.IsConnected() {
		log.Println("Failed to connect to the database")
	}

	// Initialize the HTTP server and set up routes
	http.HandleFunc("/", handler.LoginPage)
	http.HandleFunc("/login", handler.HandleLogin)
	http.HandleFunc("/register", handler.RegisterPage)
	http.HandleFunc("/register_submit", handler.HandleRegister)
	http.HandleFunc("/dashboard", handler.DashboardPage)

	// If you have files like static/style.css or static/script.js
	// then you can serve them using the following line.
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	println("Server started at http://localhost:8800")
	http.ListenAndServe(":8800", nil)
}
