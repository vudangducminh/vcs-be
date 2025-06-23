package server

import (
	"context"
	"elasticsearch/elastic"
	"elasticsearch/object"
	"fmt"
	"net/http"
	"strings"
)

func GetAccountInfoByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("")
	searchResp, err := elastic.Es.Search(
		elastic.Es.Search.WithContext(context.Background()),
		elastic.Es.Search.WithIndex("account"),
		elastic.Es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
			"query": {
				"match": {
					"_username": "%s"
				}
			}
		}`, username))),
		elastic.Es.Search.WithTrackTotalHits(true),
		elastic.Es.Search.WithPretty(),
	)
	if err != nil {
		http.Error(w, "Error searching account", http.StatusInternalServerError)
		return
	}
	defer searchResp.Body.Close()
	if searchResp.IsError() {
		http.Error(w, "Error searching account", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", searchResp.String())
}

func AddAccountInfo(account object.Account) {
	_, err := elastic.Es.Index(
		"book",
		strings.NewReader(fmt.Sprintf(`{
			"fullname": "%s",
			"email": "%s",
			"username": "%s",
			"password": "%s"
		}`, account.FullName, account.Email, account.Username, account.Password)),
		elastic.Es.Index.WithRefresh("true"),
		elastic.Es.Index.WithContext(context.Background()),
		elastic.Es.Index.WithPretty(),
	)
	if err != nil {
		fmt.Println("Error adding account:", err)
		return
	}
	fmt.Println("Account added successfully")
}
