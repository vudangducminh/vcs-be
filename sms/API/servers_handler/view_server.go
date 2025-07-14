package servers_handler

import (
	"net/http"
	posgresql_query "sms/server/database/postgresql/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      View server details
// @Description  View server details by server name substring
// @Description  Example usage:			/servers/view_server/<server_name>
// @Description  Syntax for view all servers:		/servers/view_server/?all=true
// @Accept       json
// @Produce      json
// @Param        server_name path string true "Server name substring"
// @Success      200 {object} []object.Server "List of servers"
// @Router       /servers/view_server/{server_name} [get]
func ViewServer(c *gin.Context) {
	// Implementation for viewing server details
	serverName := c.Param("server_name")
	if serverName == "?all=true" {
		serverName = ""
	}
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
