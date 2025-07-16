package query

import (
	"context"
	"fmt"
	"net/http"
	"sms/object"
	elastic "sms/server/database/elasticsearch/connector"
	"strings"
)

func AddServerInfo(server object.Server) int {
	_, err := elastic.Es.Index(
		"server",
		strings.NewReader(fmt.Sprintf(`{
			"server_id": "%s",
			"server_name": "%s",
			"status": "%s",
			"uptime": "%d",
			"created_time": "%s",
			"last_updated_time": "%s",
			"ipv4": "%s",
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
