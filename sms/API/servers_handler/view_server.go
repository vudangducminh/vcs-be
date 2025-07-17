package servers_handler

import (
	"log"
	"net/http"
	"sms/object"
	elastic_query "sms/server/database/elasticsearch/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      View server details
// @Description  View server details with optional filtering
// @Accept       json
// @Produce      json
// @Param        filter query string false "Filter by server_id, server_name, ipv4, or status"
// @Param        string path string false "Substring to search in server_id, server_name, ipv4, or status"
// @Success      200 {object} object.Server "Server details retrieved successfully"
// @Router       /servers/view_servers/{filter}/{string} [get]
func ViewServer(c *gin.Context) {
	filter := c.Query("filter")
	str := c.Param("string")
	if str == "undefined" {
		str = ""
	}
	log.Printf("Received request to view server with filter '%s' and substring: '%s'", filter, str)
	var servers []object.Server
	var httpStatus int
	switch filter {
	case "server_id":
		servers, httpStatus = elastic_query.GetServerByIdSubstr(str)
	case "server_name":
		servers, httpStatus = elastic_query.GetServerByNameSubstr(str)
	case "ipv4":
		servers, httpStatus = elastic_query.GetServerByIPv4Substr(str)
	case "status":
		servers, httpStatus = elastic_query.GetServerByStatus(str)
	default:
		servers, httpStatus = elastic_query.GetServerByNameSubstr(str)
	}
	if httpStatus == http.StatusNotFound {
		c.JSON(http.StatusOK, gin.H{"message": "No servers found with the given requirements"})
		return
	} else if httpStatus != http.StatusOK {
		c.JSON(httpStatus, gin.H{"error": "Failed to retrieve server details"})
		return
	}

	// Prepare the response
	var response []gin.H
	for _, server := range servers {
		response = append(response, gin.H{
			"server_id":         server.ServerId,
			"server_name":       server.ServerName,
			"status":            server.Status,
			"uptime":            server.Uptime,
			"created_time":      server.CreatedTime,
			"last_updated_time": server.LastUpdatedTime,
			"ipv4":              server.IPv4,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"servers": response,
		"message": "Servers retrieved successfully",
	})
}
