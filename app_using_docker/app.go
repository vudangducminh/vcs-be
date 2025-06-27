package main

import (
	"app_using_docker/elasticsearch/elastic"
	"app_using_docker/handler"
	"net/http"
)

func main() {
	// Connect to Elasticsearch
	elastic.ConnectToEs()

	if elastic.Es == nil {
		panic("Failed to connect to Elasticsearch")
	}

	// Initialize the HTTP server and set up routes
	http.HandleFunc("/", handler.LoginPage)
	http.HandleFunc("/login", handler.HandleLogin)
	http.HandleFunc("/register", handler.RegisterPage)
	http.HandleFunc("/register_submit", handler.HandleRegister)
	http.HandleFunc("/home", handler.HomePage)

	// If you have files like static/style.css or static/script.js
	// then you can serve them using the following line.
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
