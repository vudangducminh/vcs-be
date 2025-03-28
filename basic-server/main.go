package main

import (
	"basic-server/database"
	"basic-server/server"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	database.ConnectToDB()
	http.HandleFunc("/q", server.GetPlayerScoreOnServer)   // -> localhost:port/q?...
	http.HandleFunc("/add", server.AddPlayerScoreOnServer) // -> localhost:port/add?...

	fmt.Println("Server listening on port :5002")
	log.Fatal(http.ListenAndServe(":5002", nil))
}
