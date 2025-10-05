package query

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server_service/entities"
	elastic "server_service/infrastructure/elasticsearch/connector"
	"server_service/src/algorithm"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func CheckServerExists(IPv4 string) bool {
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"ipv4": "%s"
			}
		}
	}`, IPv4)

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithPretty(),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		return false
	}

	defer res.Body.Close()

	if res.IsError() {
		return false
	}

	var searchResult struct {
		Hits struct {
			Hits []struct {
				ID string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return false
	}

	return len(searchResult.Hits.Hits) > 0
}

func AddServerInfo(server entities.Server) int {
	if CheckServerExists(server.IPv4) {
		return http.StatusConflict
	}
	_, err := elastic.Es.Index(
		"server",
		strings.NewReader(fmt.Sprintf(`{
			"server_name": "%s",
			"status": "%s",
			"uptime": %d,
			"created_time": %d,
			"last_updated_time": %d,
			"ipv4": "%s"
		}`, server.ServerName, server.Status, server.Uptime,
			server.CreatedTime, server.LastUpdatedTime, server.IPv4)),
		elastic.Es.Index.WithDocumentID(algorithm.SHA256Hash(server.IPv4)),
		elastic.Es.Index.WithRefresh("true"),
		elastic.Es.Index.WithContext(context.Background()),
		elastic.Es.Index.WithPretty(),
	)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusCreated
}

func ParseSearchResults(res *esapi.Response) ([]entities.Server, int) {
	var searchResult struct {
		Hits struct {
			Hits []struct {
				Id     string `json:"_id"`
				Source struct {
					ServerName      string `json:"server_name"`
					Status          string `json:"status"`
					Uptime          []int  `json:"uptime"`
					CreatedTime     int64  `json:"created_time"`
					LastUpdatedTime int64  `json:"last_updated_time"`
					IPv4            string `json:"ipv4"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, http.StatusInternalServerError
	}
	var servers []entities.Server
	for _, hit := range searchResult.Hits.Hits {
		server := entities.Server{
			Id:              hit.Id,
			ServerName:      hit.Source.ServerName,
			Status:          hit.Source.Status,
			Uptime:          hit.Source.Uptime,
			CreatedTime:     hit.Source.CreatedTime,
			LastUpdatedTime: hit.Source.LastUpdatedTime,
			IPv4:            hit.Source.IPv4,
		}
		servers = append(servers, server)
	}
	if len(servers) == 0 {
		return nil, http.StatusNotFound
	}
	return servers, http.StatusOK
}

func GetServerByNameSubstr(substr string) ([]entities.Server, int) {
	query := fmt.Sprintf(`{
        "size": 10000,
		"query": {
			"wildcard": {
				"server_name": {
					"value": "*%s*"
				}
			}
		}
	}`, substr)

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithPretty(),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, http.StatusNotFound
	}

	return ParseSearchResults(res)
}

func GetServerByIPv4Substr(substr string) ([]entities.Server, int) {
	query := fmt.Sprintf(`{
        "size": 10000,
		"query": {
			"wildcard": {
				"ipv4": {
					"value": "*%s*"
				}
			}
		}
	}`, substr)

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithPretty(),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, http.StatusNotFound
	}

	return ParseSearchResults(res)
}

func GetServerByStatus(substr string) ([]entities.Server, int) {
	query := fmt.Sprintf(`{
        "size": 10000,
		"query": {
			"match": {
				"status": "%s"
			}
		}
	}`, substr)

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithPretty(),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		return nil, http.StatusInternalServerError
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, http.StatusNotFound
	}

	return ParseSearchResults(res)
}

func GetServerById(Id string) (entities.Server, bool) {
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"_id": "%s"
			}
		}
	}`, Id)

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithPretty(),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		return entities.Server{}, false
	}

	defer res.Body.Close()

	if res.IsError() {
		return entities.Server{}, false
	}

	results, status := ParseSearchResults(res)
	if status != http.StatusOK || len(results) == 0 {
		return entities.Server{}, false
	}

	return results[0], true
}

func UpdateServerInfo(server entities.Server) int {
	uptimeJSON, _ := json.Marshal(server.Uptime)
	_, err := elastic.Es.Update(
		"server",
		server.Id,
		strings.NewReader(fmt.Sprintf(`{
			"doc": {
				"server_name": "%s",
				"status": "%s",
				"uptime": %s,
				"last_updated_time": %d,
				"ipv4": "%s"
			}
		}`, server.ServerName, server.Status, string(uptimeJSON), server.LastUpdatedTime, server.IPv4)),
		elastic.Es.Update.WithContext(context.Background()),
		elastic.Es.Update.WithPretty(),
	)

	if err != nil {
		log.Println("Update error:", err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func DeleteServer(Id string) int {
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"_id": "%s"
			}
		}
	}`, Id)

	res, err := elastic.Es.DeleteByQuery(
		[]string{"server"},
		strings.NewReader(query),
		elastic.Es.DeleteByQuery.WithContext(context.Background()),
		elastic.Es.DeleteByQuery.WithRefresh(true),
	)

	if err != nil {
		return http.StatusInternalServerError
	}
	defer res.Body.Close()

	if res.IsError() {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func BulkServerInfo(servers []entities.Server) int {
	var bulkRequest strings.Builder
	for _, server := range servers {
		bulkRequest.WriteString(fmt.Sprintf(`{"index": {"_index": "server", "_id": "%s"}}%s`, algorithm.SHA256Hash(server.IPv4), "\n"))
		uptimeJSON, _ := json.Marshal(server.Uptime)
		bulkRequest.WriteString(fmt.Sprintf(`{"server_name": "%s", "status": "%s", "uptime": %s, "created_time": %d, "last_updated_time": %d, "ipv4": "%s"}%s`,
			server.ServerName, server.Status, string(uptimeJSON), server.CreatedTime, server.LastUpdatedTime, server.IPv4, "\n"))
	}

	if len(bulkRequest.String()) == 0 {
		log.Println("No servers to index")
		return http.StatusCreated
	}
	res, err := elastic.Es.Bulk(
		strings.NewReader(bulkRequest.String()),
		elastic.Es.Bulk.WithIndex("server"),
		elastic.Es.Bulk.WithContext(context.Background()),
		elastic.Es.Bulk.WithPretty(),
	)
	// log.Println("Response:", res)
	// log.Println("Error:", err)
	if err != nil {
		return http.StatusInternalServerError
	}
	defer res.Body.Close()

	if res.IsError() {
		return http.StatusInternalServerError
	}

	return http.StatusCreated
}
