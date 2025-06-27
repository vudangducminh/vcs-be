package query

import (
	"app_using_docker/elasticsearch/elastic"
	"app_using_docker/object"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
					"username": "%s"
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

func GetAccountPasswordByUsername(username string) string {
	searchResp, err := elastic.Es.Search(
		elastic.Es.Search.WithContext(context.Background()),
		elastic.Es.Search.WithIndex("account"),
		elastic.Es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
			"query": {
				"match": {
					"username": "%s"
				}
			}
		}`, username))),
		elastic.Es.Search.WithTrackTotalHits(true),
		elastic.Es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalln("Error searching account")
		return ""
	}
	defer searchResp.Body.Close()
	if searchResp.IsError() {
		log.Fatalln("Error searching account")
		return ""
	}
	var searchResult map[string]interface{}
	err = json.NewDecoder(searchResp.Body).Decode(&searchResult)
	if err != nil {
		log.Fatalln("Error decoding search response")
		return ""
	}
	hits, ok := searchResult["hits"].(map[string]interface{})
	if !ok {
		log.Fatalln("Invalid search response format")
		return ""
	}
	newHits, ok := hits["hits"].([]interface{})
	if !ok {
		log.Fatalln("Invalid search response format")
		return ""
	}
	if len(newHits) == 0 {
		log.Fatalln("No account found with the given username")
		return ""
	}
	firstHit, ok := newHits[0].(map[string]interface{})
	if !ok {
		log.Fatalln("Invalid search response format")
		return ""
	}
	source, ok := firstHit["_source"].(map[string]interface{})
	if !ok {
		log.Fatalln("Invalid search response format")
		return ""
	}
	password, ok := source["password"].(string)
	if !ok {
		log.Fatalln("Invalid search response format")
		return ""
	}
	return password
}

func CheckAccountExistsByUsername(username string) bool {
	value, err := elastic.Es.Search(
		elastic.Es.Search.WithContext(context.Background()),
		elastic.Es.Search.WithIndex("account"),
		elastic.Es.Search.WithBody(strings.NewReader(fmt.Sprintf(`{
			"query": {
				"match": {
					"username": "%s"
				}
			}
		}`, username))),
		elastic.Es.Search.WithTrackTotalHits(true),
		elastic.Es.Search.WithPretty(),
	)
	if err != nil {
		log.Println("Error searching account:", err)
		return true
	}
	defer value.Body.Close()
	if value.IsError() {
		log.Println("Error searching account:", err)
		return true
	}
	var searchResult map[string]interface{}
	err = json.NewDecoder(value.Body).Decode(&searchResult)
	if err != nil {
		log.Println("Error decoding search response:", err)
		return true
	}
	hits, ok := searchResult["hits"].(map[string]interface{})
	if !ok {
		log.Println("Invalid search response format")
		return true
	}
	total, ok := hits["total"].(map[string]interface{})
	if !ok {
		log.Println("Invalid search response format")
		return true
	}
	val, ok := total["value"].(float64)
	if !ok {
		log.Println("Invalid search response format")
		return true
	}
	return val > 0
}

func CountAccounts() int {
	count, err := elastic.Es.Count(
		elastic.Es.Count.WithContext(context.Background()),
		elastic.Es.Count.WithIndex("account"),
		elastic.Es.Count.WithPretty(),
	)
	if err != nil {
		log.Fatalln("Error counting accounts:", err)
		return 0
	}
	defer count.Body.Close()
	if count.IsError() {
		log.Fatalln("Error counting accounts:", err)
		return 0
	}
	var countResult map[string]interface{}
	err = json.NewDecoder(count.Body).Decode(&countResult)
	if err != nil {
		log.Fatalln("Error decoding 'count' response: ", err)
		return 0
	}
	countValue, ok := countResult["count"].(float64)
	if !ok {
		log.Fatalln("Invalid 'count' response: ", err)
		return 0
	}
	return int(countValue)
}

func AddAccountInfo(account object.Account) {
	if CheckAccountExistsByUsername(account.Username) {
		log.Println("Account with this username already exists")
		return
	}
	_, err := elastic.Es.Index(
		"account",
		strings.NewReader(fmt.Sprintf(`{
			"id": "%d",
			"fullname": "%s",
			"email": "%s",
			"username": "%s",
			"password": "%s"
		}`, CountAccounts(), account.FullName, account.Email, account.Username, account.Password)),
		elastic.Es.Index.WithRefresh("true"),
		elastic.Es.Index.WithContext(context.Background()),
		elastic.Es.Index.WithPretty(),
	)
	if err != nil {
		log.Println("Error adding account:", err)
		return
	}
	log.Println("Account added successfully")
}
