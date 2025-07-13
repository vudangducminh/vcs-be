package servers_handler

import (
	"log"
	"net/http"
	posgresql_query "sms/server/database/postgresql/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      View server details
// @Description  View server details by server name substring
// @Accept       json
// @Produce      json
// @Param        server_name path string false "Server name substring"
// @Success      200 {object} []object.Server "List of servers"
// @Router       /servers/view_server/{server_name} [get]
func ViewServer(c *gin.Context) {
	serverName := c.Param("server_name")
	if serverName == "undefined" {
		serverName = ""
	}
	log.Println("Received request to view server with name substring:", serverName)
	servers, httpStatus := posgresql_query.GetServerBySubstr(serverName)
	if httpStatus == http.StatusNotFound {
		c.JSON(http.StatusOK, gin.H{"message": "No servers found with the given name"})
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
