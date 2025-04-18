package main

import (
	"elasticsearch/elastic"
	"elasticsearch/object"
	"elasticsearch/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Elasticsearch API!")
	})
	http.HandleFunc("/book/id", server.GetBookInfoByID)
	http.HandleFunc("/book/title", server.GetBookInfoByTitle)
	http.HandleFunc("/book/search/title", server.SearchBookInfoByTitlePrefix)
	// http.HandleFunc("/add", server.AddPlayerScoreOnServer)
	elastic.ConnectToEs()
	book := object.Book{
		Title:          "Prince of Persia",
		Published_date: "2023-10-01",
		Rating:         "4.4",
		Author:         "John Doe",
	}
	server.AddBookInfo(book)
	fmt.Println("Server listening on port :5002")
	log.Fatal(http.ListenAndServe(":5002", nil))
}
