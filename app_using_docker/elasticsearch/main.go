package main

import (
	"elasticsearch/elastic"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Elasticsearch API!")
	})

	elastic.ConnectToEs()
	fmt.Println("Server listening on port: 5002")
	log.Fatal(http.ListenAndServe(":5002", nil))
}
