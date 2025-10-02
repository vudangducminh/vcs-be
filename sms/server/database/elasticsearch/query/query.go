package query

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sms/algorithm"
	"sms/object"
	elastic "sms/server/database/elasticsearch/connector"
	"sort"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func GetAllServer() []object.BriefServerInfo {
	// Check if Elasticsearch is connected before making the query
	if !elastic.IsConnected() {
		log.Println("Elasticsearch not connected, attempting to reconnect...")
		if err := elastic.Reconnect(); err != nil {
			log.Println("Failed to reconnect to Elasticsearch:", err)
			return nil
		}
		log.Println("Successfully reconnected to Elasticsearch")
	}

	query := `{
		"size": 10000,
		"_source": ["ipv4", "uptime"],
		"query": {
			"match_all": { }
		}
	}`

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithPretty(),
		elastic.Es.Search.WithContext(context.Background()),
	)

	if err != nil {
		log.Println("Error getting all servers:", err)
		return nil
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Println("Error response from Elasticsearch:", res.String())
		return nil
	}

	var searchResult struct {
		Hits struct {
			Hits []struct {
				Id     string `json:"_id"`
				Source struct {
					IPv4   string `json:"ipv4"`
					Uptime []int  `json:"uptime"`
				} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		log.Println("Error decoding search results:", err)
		return nil
	}

	var servers []object.BriefServerInfo
	for _, hit := range searchResult.Hits.Hits {
		servers = append(servers, object.BriefServerInfo{
			Id:     hit.Id,
			IPv4:   hit.Source.IPv4,
			Uptime: hit.Source.Uptime,
		})
	}

	return servers
}

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

func ParseSearchResults(res *esapi.Response) ([]object.Server, int) {
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
	var servers []object.Server
	for _, hit := range searchResult.Hits.Hits {
		server := object.Server{
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

func GetServerByIdSubstr(substr string) ([]object.Server, int) {
	query := fmt.Sprintf(`{
        "size": 10000,
		"query": {
			"wildcard": {
				"_id": {
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

func GetServerByIPv4Substr(substr string) ([]object.Server, int) {
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

func GetServerByStatus(substr string) ([]object.Server, int) {
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

func GetServerById(Id string) (object.Server, bool) {
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
		server.Id,
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

func GetTotalActiveServersCount(filter string, substr string) int {
	var query string
	switch filter {
	case "server_name":
		query = `{
			"query": {
				"bool": {
					"must": [
						{
							"wildcard": {
								"server_name": {
									"value": "*%s*"
								}
							}
						},
						{
							"term": {
								"status": "active"
							}
						}
					]
				}
			}
		}`
	case "ipv4":
		query = `{
			"query": {
				"bool": {
					"must": [
						{
							"wildcard": {
								"ipv4": {
									"value": "*%s*"
								}
							}
						},
						{
							"term": {
								"status": "active"
							}
						}
					]
				}
			}
		}`
	default:
		query = `{
			"query": {
				"term": {
					"status": "active"
				}
			}
		}`
	}

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

func GetTotalInactiveServersCount(filter string, substr string) int {
	var query string
	switch filter {
	case "server_name":
		query = `{
			"query": {
				"bool": {
					"must": [
						{
							"wildcard": {
								"server_name": {
									"value": "*%s*"
								}
							}
						},
						{
							"term": {
								"status": "inactive"
							}
						}
					]
				}
			}
		}`
	case "ipv4":
		query = `{
			"query": {
				"bool": {
					"must": [
						{
							"wildcard": {
								"ipv4": {
									"value": "*%s*"
								}
							}
						},
						{
							"term": {
								"status": "inactive"
							}
						}
					]
				}
			}
		}`
	default:
		query = `{
			"query": {
				"term": {
					"status": "inactive"
				}
			}
		}`
	}

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

func GetTotalMaintenanceServersCount(filter string, substr string) int {
	var query string
	switch filter {
	case "server_name":
		query = `{
			"query": {
				"bool": {
					"must": [
						{
							"wildcard": {
								"server_name": {
									"value": "*%s*"
								}
							}
						},
						{
							"term": {
								"status": "maintenance"
							}
						}
					]
				}
			}
		}`
	case "ipv4":
		query = `{
			"query": {
				"bool": {
					"must": [
						{
							"wildcard": {
								"ipv4": {
									"value": "*%s*"
								}
							}
						},
						{
							"term": {
								"status": "maintenance"
							}
						}
					]
				}
			}
		}`
	default:
		query = `{
			"query": {
				"term": {
					"status": "maintenance"
				}
			}
		}`
	}

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
	log.Println("Getting total last updated time...")
	query := `{
		"size": 0,
		"query": {
			"term": {
				"status": "active"
			}
		},
		"aggs": {
			"total_last_updated_time": {
				"sum": {
					"field": "last_updated_time"
				}
			}
		}
	}`

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithContext(context.Background()),
		elastic.Es.Search.WithPretty(),
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

func BulkServerInfo(servers []object.Server) int {
	var bulkRequest strings.Builder
	for _, server := range servers {
		bulkRequest.WriteString(fmt.Sprintf(`{"index": {"_index": "server", "_id": "%s"}}%s`, algorithm.SHA256Hash(server.IPv4), "\n"))
		bulkRequest.WriteString(fmt.Sprintf(`{"server_name": "%s", "status": "%s", "uptime": %d, "created_time": %d, "last_updated_time": %d, "ipv4": "%s"}%s`,
			server.ServerName, server.Status, server.Uptime, server.CreatedTime, server.LastUpdatedTime, server.IPv4, "\n"))
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

func BulkUpdateServerInfo(updates []object.ServerUptimeUpdate) int {
	var bulkRequest strings.Builder

	for _, update := range updates {
		// Update action
		bulkRequest.WriteString(fmt.Sprintf(`{"update": {"_index": "server", "_id": "%s"}}%s`, update.Id, "\n"))
		uptimeJSON, _ := json.Marshal(update.Uptime)
		// Document to update
		bulkRequest.WriteString(fmt.Sprintf(`{"doc": {"uptime": %s, "last_updated_time": %d, "status": "%s"}}%s`,
			string(uptimeJSON), time.Now().Unix(), update.Status, "\n"))
	}

	if len(bulkRequest.String()) == 0 {
		log.Println("No updates to process")
		return http.StatusOK
	}

	res, err := elastic.Es.Bulk(
		strings.NewReader(bulkRequest.String()),
		elastic.Es.Bulk.WithIndex("server"),
		elastic.Es.Bulk.WithContext(context.Background()),
		elastic.Es.Bulk.WithPretty(),
	)

	if err != nil {
		log.Println("Bulk update error:", err)
		return http.StatusInternalServerError
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Println("Bulk update response error:", res.String())
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func GetServerUptimeInRange(startBlock int, endBlock int, order string, filter string, substr string) ([]object.Server, int, float64) {
	var query string
	switch filter {
	case "server_name":
		query = `{
			"size": 10000,
			"query": {
				"wildcard": {
					"server_name": {
						"value": "*` + substr + `*"
					}
				}
			}
		}`
	case "ipv4":
		query = `{
			"size": 10000,
			"query": {
				"wildcard": {
					"ipv4": {
						"value": "*` + substr + `*"
					}
				}
			}
		}`
	case "status":
		query = `{
			"size": 10000,
			"query": {
				"match": {
					"status": "` + substr + `"
				}
			}
		}`
	default:
		query = `{
			"size": 10000,
			"query": {
				"match_all": { }
			}
		}`
	}

	res, err := elastic.Es.Search(
		elastic.Es.Search.WithIndex("server"),
		elastic.Es.Search.WithBody(strings.NewReader(query)),
		elastic.Es.Search.WithPretty(),
		elastic.Es.Search.WithContext(context.Background()),
	)
	if err != nil {
		return nil, http.StatusInternalServerError, 0
	}

	defer res.Body.Close()
	if res.IsError() {
		return nil, http.StatusNotFound, 0
	}

	servers, status := ParseSearchResults(res)
	if status != http.StatusOK {
		return nil, status, 0
	}
	// log.Println("Start block: ", startBlock)
	// log.Println("End block: ", endBlock)
	var allServerTotalTime float64 = 0
	var allServerUptime float64 = 0
	for i := 0; i < len(servers); i++ {
		// log.Println("Server IP: ", servers[i].IPv4)
		// log.Println("Uptime data: ", servers[i].Uptime)
		var start = max(0, len(servers[i].Uptime)-startBlock)
		var end = max(0, len(servers[i].Uptime)-endBlock)
		allServerTotalTime += float64((end - start + 1) * 1200)
		// Calculate total uptime in the range using simple loop
		var totalUptime int = 0
		for j := start; j <= end && j < len(servers[i].Uptime); j++ {
			totalUptime += servers[i].Uptime[j]
		}
		allServerUptime += float64(totalUptime)
		// Update the server with calculated uptime
		servers[i].Uptime = []int{totalUptime} // Convert back to expected format
	}

	sort.Slice(servers, func(i, j int) bool {
		var less bool
		switch filter {
		case "status":
			less = servers[i].Status < servers[j].Status
		case "ipv4":
			less = servers[i].IPv4 < servers[j].IPv4
		default: // Default to sorting by server_name
			less = servers[i].ServerName < servers[j].ServerName
		}
		if order == "desc" {
			return !less
		}
		return less
	})

	log.Println("All server total time:", allServerTotalTime)
	log.Println("All server total uptime:", allServerUptime)
	log.Println("Uptime percentage:", allServerUptime/allServerTotalTime*100)
	return servers, http.StatusOK, allServerUptime / allServerTotalTime * 100
}
