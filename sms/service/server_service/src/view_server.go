package src

import (
	"log"
	"net/http"
	"server_service/entities"
	elastic_query "server_service/infrastructure/elasticsearch/query"
	"sort"

	"github.com/gin-gonic/gin"
)

// @Tags         Server
// @Summary      View server details
// @Description  View server details with optional filtering
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        order query string false "Order of results, either 'asc' or 'desc'. If not provided or using the wrong order format, the default order is ascending"
// @Param        filter query string false "Filter by server_name, ipv4, or status. If not provided or using the wrong filter format, the default filter is server_name"
// @Param        string query string false "Substring to search in server_name, ipv4, or status"
// @Success      200 {object} entities.ViewServerSuccessResponse "Server details retrieved successfully"
// @Failure      400 {object} entities.ViewServerBadRequestResponse "Invalid request parameters"
// @Failure      401 {object} entities.AuthErrorResponse "Authentication failed"
// @Failure      429 {object} entities.RateLimitExceededResponse "Too many requests"
// @Failure      500 {object} entities.ViewServerInternalServerErrorResponse "Failed to retrieve server details"
// @Router       /servers/view-servers [get]
func ViewServer(c *gin.Context) {

	order := c.Query("order")
	if order != "asc" && order != "desc" {
		order = "asc" // Default order if not specified
	}
	filter := c.Query("filter")
	str := c.Query("string")
	log.Printf("Received request to view server with filter '%s' and substring: '%s'", filter, str)
	if str == "undefined" || str == "{string}" {
		str = ""
	}
	log.Printf("Received request to view server with filter '%s' and substring: '%s'", filter, str)
	var servers []entities.Server
	var httpStatus int
	switch filter {
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
			"_id":               server.Id,
			"server_name":       server.ServerName,
			"status":            server.Status,
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

type ServerViewer interface {
	GetServerByNameSubstr(str string) ([]entities.Server, int)
	GetServerByIPv4Substr(str string) ([]entities.Server, int)
	GetServerByStatus(str string) ([]entities.Server, int)
}

var serverViewer ServerViewer

func SetServerViewer(sv ServerViewer) {
	serverViewer = sv
}

func ModifiedViewServer(c *gin.Context) {
	order := c.Query("order")
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	filter := c.Query("filter")
	str := c.Query("string")
	log.Printf("Received request to view server with filter '%s' and substring: '%s'", filter, str)
	if str == "undefined" || str == "{string}" {
		str = ""
	}

	var servers []entities.Server
	var httpStatus int
	switch filter {
	case "server_name":
		servers, httpStatus = serverViewer.GetServerByNameSubstr(str)
	case "ipv4":
		servers, httpStatus = serverViewer.GetServerByIPv4Substr(str)
	case "status":
		servers, httpStatus = serverViewer.GetServerByStatus(str)
	default:
		servers, httpStatus = serverViewer.GetServerByNameSubstr(str)
	}
	if httpStatus == http.StatusNotFound {
		c.JSON(http.StatusOK, gin.H{"message": "No servers found with the given requirements"})
		return
	} else if httpStatus != http.StatusOK {
		c.JSON(httpStatus, gin.H{"error": "Failed to retrieve server details"})
		return
	}

	sort.Slice(servers, func(i, j int) bool {
		var less bool
		switch filter {
		case "status":
			less = servers[i].Status < servers[j].Status
		case "ipv4":
			less = servers[i].IPv4 < servers[j].IPv4
		default:
			less = servers[i].ServerName < servers[j].ServerName
		}
		if order == "desc" {
			return !less
		}
		return less
	})

	var response []gin.H
	for _, server := range servers {
		response = append(response, gin.H{
			"_id":               server.Id,
			"server_name":       server.ServerName,
			"status":            server.Status,
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
