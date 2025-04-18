package server

import (
	"context"
	"elasticsearch/elastic"
	"elasticsearch/object"
	"fmt"
	"net/http"
	"strconv"
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
				"match": {
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

func SearchRelatedBookInfoByTitle(w http.ResponseWriter, r *http.Request) {
	titlePart := r.URL.Query().Get("")
	if titlePart == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}
	searchResp, err := elastic.Es.Search(
		elastic.Es.Search.WithContext(context.Background()),
		elastic.Es.Search.WithIndex("book"),
		elastic.Es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
			"query": {
				"wildcard": {
					"title": {
						"value": "*%s*",
						"case_insensitive": true
					}
				}
			}
		}`, titlePart))),
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

func SearchBookInfoByRating(w http.ResponseWriter, r *http.Request) {
	rating := r.URL.Query().Get("")
	if rating == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}
	ratingFloat, err := strconv.ParseFloat(rating, 64)
	if err != nil {
		http.Error(w, "Invalid rating parameter", http.StatusBadRequest)
		return
	}

	searchResp, err := elastic.Es.Search(
		elastic.Es.Search.WithContext(context.Background()),
		elastic.Es.Search.WithIndex("book"),
		elastic.Es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
			"query": {
				"range": {
					"rating": {
						"gte": "%f"
					}
				}
			}
		}`, ratingFloat))),
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

func DeleteBookInfoByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	deleteResp, err := elastic.Es.Delete(
		"book",
		id,
		elastic.Es.Delete.WithRefresh("true"),
		elastic.Es.Delete.WithPretty(),
		elastic.Es.Delete.WithContext(context.Background()),
	)
	if err != nil {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
		return
	}
	defer deleteResp.Body.Close()
	if deleteResp.IsError() {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", deleteResp.String())
}

func AddBookInfo(book object.Book) {
	_, err := elastic.Es.Index(
		"book",
		strings.NewReader(fmt.Sprintf(`{
			"title": "%s",
			"published_date": "%s",
			"rating": "%s",
			"author": "%s"
		}`, book.Title, book.Published_date, book.Rating, book.Author)),
		elastic.Es.Index.WithRefresh("true"),
		elastic.Es.Index.WithContext(context.Background()),
		elastic.Es.Index.WithPretty(),
	)
	if err != nil {
		fmt.Println("Error adding book:", err)
		return
	}
	fmt.Println("Book added successfully")
}
