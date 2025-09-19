package servers_handler

import (
	"log"
	"net/http"
	"sms/object"
	elastic_query "sms/server/database/elasticsearch/query"
	"sort"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      View server details
// @Description  View server details with optional filtering
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        order query string false "Order of results, either 'asc' or 'desc'. If not provided or using the wrong order format, the default order is ascending"
// @Param        filter query string false "Filter by server_id, server_name, ipv4, or status. If not provided or using the wrong filter format, the default filter is server_name"
// @Param        string query string false "Substring to search in server_id, server_name, ipv4, or status"
// @Success      200 {object} object.ViewServerSuccessResponse "Server details retrieved successfully"
// @Failure      400 {object} object.ViewServerBadRequestResponse "Invalid request parameters"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      500 {object} object.ViewServerInternalServerErrorResponse "Failed to retrieve server details"
// @Router       /servers/view_servers/{order}/{filter}/{string} [get]
func ViewServer(c *gin.Context) {

	order := c.Query("order")
	if order != "asc" && order != "desc" {
		order = "asc" // Default order if not specified
	}
	filter := c.Query("filter")
	str := c.Param("string")
	log.Printf("Received request to view server with filter '%s' and substring: '%s'", filter, str)
	if str == "undefined" || str == "{string}" {
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

	// Sort the servers based on the filter and order
	sort.Slice(servers, func(i, j int) bool {
		var less bool
		switch filter {
		case "server_id":
			less = servers[i].ServerId < servers[j].ServerId
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
