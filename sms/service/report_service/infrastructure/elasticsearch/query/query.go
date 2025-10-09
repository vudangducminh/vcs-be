package query

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"report_service/entities"
	elastic "report_service/infrastructure/elasticsearch/connector"
	"sort"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

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

func GetServerUptimeInRange(startBlock int, endBlock int, order string, filter string, substr string) ([]entities.Server, int, float64) {
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
	case "status":
		query = `{
			"query": {
				"term": {
					"status": "%s"
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
	case "status":
		query = `{
			"query": {
				"term": {
					"status": "%s"
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
	case "status":
		query = `{
			"query": {
				"term": {
					"status": "%s"
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
