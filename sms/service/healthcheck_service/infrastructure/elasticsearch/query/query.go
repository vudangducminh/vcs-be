package query

import (
	"context"
	"encoding/json"
	"fmt"
	"healthcheck_service/entities"
	elastic "healthcheck_service/infrastructure/elasticsearch/connector"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetAllServer() []entities.BriefServerInfo {
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

	var servers []entities.BriefServerInfo
	for _, hit := range searchResult.Hits.Hits {
		servers = append(servers, entities.BriefServerInfo{
			Id:     hit.Id,
			IPv4:   hit.Source.IPv4,
			Uptime: hit.Source.Uptime,
		})
	}

	return servers
}

func BulkUpdateServerInfo(updates []entities.ServerUptimeUpdate) int {
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
