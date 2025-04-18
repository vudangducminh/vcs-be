package main

import (
	"elasticsearch/elastic"
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
	http.HandleFunc("/book/search/title", server.SearchRelatedBookInfoByTitle)
	http.HandleFunc("/book/search/rating", server.SearchBookInfoByRating)
	http.HandleFunc("/book/delete/id", server.DeleteBookInfoByID)

	elastic.ConnectToEs()

	// book := object.Book{
	// 	Title:          "[MV] Mumei",
	// 	Published_date: "2023-10-10",
	// 	Rating:         "9.6",
	// 	Author:         "Nanashi Mumei",
	// }
	// server.AddBookInfo(book)
	fmt.Println("Server listening on port: 5002")
	log.Fatal(http.ListenAndServe(":5002", nil))
}
