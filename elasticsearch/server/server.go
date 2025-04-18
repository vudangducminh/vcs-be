package server

import (
	"context"
	"elasticsearch/elastic"
	"elasticsearch/object"
	"fmt"
	"net/http"
	"strings"
)

func GetBookInfoByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	searchResp, err := elastic.Es.Search(
		elastic.Es.Search.WithContext(context.Background()),
		elastic.Es.Search.WithIndex("book"),
		elastic.Es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
			"query": {
				"term": {
					"_id": "%s"
				}
			}
		}`, id))),
		elastic.Es.Search.WithTrackTotalHits(true),
		elastic.Es.Search.WithPretty(),
	)
	if err != nil {
		http.Error(w, "Error searching for book", http.StatusInternalServerError)
		return
	}
	defer searchResp.Body.Close()
	if searchResp.IsError() {
		http.Error(w, "Error searching for book", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", searchResp.String())
}

func GetBookInfoByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("")
	if title == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}
	searchResp, err := elastic.Es.Search(
		elastic.Es.Search.WithContext(context.Background()),
		elastic.Es.Search.WithIndex("book"),
		elastic.Es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
			"query": {
				"match": {
					"title": "%s"
				}
		}
		}`, title))),
		elastic.Es.Search.WithTrackTotalHits(true),
		elastic.Es.Search.WithPretty(),
	)
	if err != nil {
		http.Error(w, "Error searching for book", http.StatusInternalServerError)
		return
	}
	defer searchResp.Body.Close()
	if searchResp.IsError() {
		http.Error(w, "Error searching for book", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", searchResp.String())
}

func SearchBookInfoByTitlePrefix(w http.ResponseWriter, r *http.Request) {
	titlePrefix := r.URL.Query().Get("")
	titlePrefix = strings.ToLower(titlePrefix)
	if titlePrefix == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}
	searchResp, err := elastic.Es.Search(
		elastic.Es.Search.WithContext(context.Background()),
		elastic.Es.Search.WithIndex("book"),
		elastic.Es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
			"query": {
				"prefix": {
					"title": "%s"
				}
		}
		}`, titlePrefix))),
		elastic.Es.Search.WithTrackTotalHits(true),
		elastic.Es.Search.WithPretty(),
	)
	if err != nil {
		http.Error(w, "Error searching for book", http.StatusInternalServerError)
		return
	}
	defer searchResp.Body.Close()
	if searchResp.IsError() {
		http.Error(w, "Error searching for book", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", searchResp.String())
}

func AddBookInfo(book object.Book) {
	book.Title = strings.ToLower(book.Title)
	book.Author = strings.ToLower(book.Author)
	_, err := elastic.Es.Index(
		"book",
		strings.NewReader(fmt.Sprintf(`{
			"title": "%s",
			"published_date": "%s",
			"rating": "%s",
			"author": "%s"
		}`, book.Title, book.Published_date, book.Rating, book.Author)),
		elastic.Es.Index.WithRefresh("true"),
	)
	if err != nil {
		fmt.Println("Error adding book:", err)
		return
	}
	fmt.Println("Book added successfully")
}
