package query

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sms/object"
	elastic "sms/server/database/elasticsearch/connector"
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

func AddServerInfo(server object.Server) int {
	if CheckServerExists(server.IPv4) {
		return http.StatusConflict
	}
	_, err := elastic.Es.Index(
		"server",
		strings.NewReader(fmt.Sprintf(`{
			"server_id": "%s",
			"server_name": "%s",
			"status": "%s",
			"uptime": %d,
			"created_time": %d,
			"last_updated_time": %d,
			"ipv4": "%s"
		}`, server.ServerId, server.ServerName, server.Status, server.Uptime,
			server.CreatedTime, server.LastUpdatedTime, server.IPv4)),
		elastic.Es.Index.WithRefresh("true"),
		elastic.Es.Index.WithContext(context.Background()),
		elastic.Es.Index.WithPretty(),
	)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusCreated
}

func ParseSearchResults(res *esapi.Response) ([]object.Server, int) {
	var searchResult struct {
		Hits struct {
			Hits []struct {
				Source struct {
					ServerId        string `json:"server_id"`
					ServerName      string `json:"server_name"`
					Status          string `json:"status"`
					Uptime          int64  `json:"uptime"`
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
	var servers []object.Server
	for _, hit := range searchResult.Hits.Hits {
		server := object.Server{
			ServerId:        hit.Source.ServerId,
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

func GetServerByIdSubstr(substr string) ([]object.Server, int) {
	query := fmt.Sprintf(`{
		"query": {
			"wildcard": {
				"server_id": {
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

func GetServerByNameSubstr(substr string) ([]object.Server, int) {
	query := fmt.Sprintf(`{
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

func GetServerByIPv4Substr(substr string) ([]object.Server, int) {
	query := fmt.Sprintf(`{
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

func GetServerByStatus(substr string) ([]object.Server, int) {
	query := fmt.Sprintf(`{
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

func GetServerById(serverId string) (object.Server, bool) {
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"server_id": "%s"
			}
		}
	}`, serverId)

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithPretty(),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		return object.Server{}, false
	}

	defer res.Body.Close()

	if res.IsError() {
		return object.Server{}, false
	}

	results, status := ParseSearchResults(res)
	if status != http.StatusOK || len(results) == 0 {
		return object.Server{}, false
	}

	return results[0], true
}

func UpdateServerInfo(server object.Server) int {
	_, err := elastic.Es.Update(
		"server",
		server.ServerId,
		strings.NewReader(fmt.Sprintf(`{
			"doc": {
				"server_name": "%s",
				"status": "%s",
				"uptime": %d,
				"last_updated_time": %d,
				"ipv4": "%s"
			}
		}`, server.ServerName, server.Status, server.Uptime, server.LastUpdatedTime, server.IPv4)),
		elastic.Es.Update.WithContext(context.Background()),
		elastic.Es.Update.WithPretty(),
	)

	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func DeleteServer(serverId string) int {
	query := fmt.Sprintf(`{
		"query": {
			"term": {
				"server_id": "%s"
			}
		}
	}`, serverId)

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

func GetTotalServersCount() int {
	res, err := elastic.Es.Count(
		elastic.Es.Count.WithIndex("server"),
		elastic.Es.Count.WithContext(context.Background()),
	)

	if err != nil {
		return 0
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0
	}

	var countResult struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0
	}

	return countResult.Count
}

func GetTotalActiveServersCount() int {
	query := `{
		"query": {
			"term": {
				"status": "active"
			}
		}
	}`

	res, err := elastic.Es.Count(
		elastic.Es.Count.WithIndex("server"),
		elastic.Es.Count.WithBody(strings.NewReader(query)),
		elastic.Es.Count.WithContext(context.Background()),
	)

	if err != nil {
		return 0
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0
	}

	var countResult struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0
	}

	return countResult.Count
}

func GetTotalInactiveServersCount() int {
	query := `{
		"query": {
			"term": {
				"status": "inactive"
			}
		}
	}`

	res, err := elastic.Es.Count(
		elastic.Es.Count.WithIndex("server"),
		elastic.Es.Count.WithBody(strings.NewReader(query)),
		elastic.Es.Count.WithContext(context.Background()),
	)

	if err != nil {
		return 0
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0
	}

	var countResult struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0
	}

	return countResult.Count
}

func GetTotalMaintenanceServersCount() int {
	query := `{
		"query": {
			"term": {
				"status": "maintenance"
			}
		}
	}`

	res, err := elastic.Es.Count(
		elastic.Es.Count.WithIndex("server"),
		elastic.Es.Count.WithBody(strings.NewReader(query)),
		elastic.Es.Count.WithContext(context.Background()),
	)

	if err != nil {
		return 0
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0
	}

	var countResult struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0
	}

	return countResult.Count
}

func GetTotalCreatedTime() (int64, int) {
	query := `{
		"size": 0,
		"aggs": {
			"total_created_time": {
				"sum": {
					"field": "created_time"
				}
			}
		}
	}`

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		return 0, 0
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, 0
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
		Aggregations struct {
			TotalCreatedTime struct {
				Value float64 `json:"value"`
			} `json:"total_created_time"`
		} `json:"aggregations"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, 0
	}

	return int64(result.Aggregations.TotalCreatedTime.Value), result.Hits.Total.Value
}

func GetTotalLastUpdatedTime() (int64, int) {
	query := `{
		"size": 0,
		"query": {
			"term": {
				"status": "active"
			}
		},
		"aggs": {
			"total_created_time": {
				"sum": {
					"field": "uptime"
				}
			},
		}
	}`

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		return 0, 0
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, 0
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
		Aggregations struct {
			TotalLastUpdatedTime struct {
				Value float64 `json:"value"`
			} `json:"total_last_updated_time"`
		} `json:"aggregations"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, 0
	}

	return int64(result.Aggregations.TotalLastUpdatedTime.Value), result.Hits.Total.Value
}

func GetTotalUptime() (int64, int) {
	query := `{
		"size": 0,
		"aggs": {
			"total_uptime": {
				"sum": {
					"field": "uptime"
				}
			}
		}
	}`

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		return 0, 0
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, 0
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
		} `json:"hits"`
		Aggregations struct {
			TotalUptime struct {
				Value float64 `json:"value"`
			} `json:"total_uptime"`
		} `json:"aggregations"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, 0
	}

	return int64(result.Aggregations.TotalUptime.Value), result.Hits.Total.Value
}
