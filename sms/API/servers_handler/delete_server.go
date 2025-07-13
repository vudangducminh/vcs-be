package servers_handler

import (
	"net/http"

	posgresql_query "sms/server/database/postgresql/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      Delete a server by ID
// @Description  Delete a server by its unique ID
// @Accept       json
// @Produce      json
// @Param        server_id path string true "Server ID"
// @Success      200 {object} object.DeleteServerResponse "Server deleted successfully"
// @Router	     /servers/delete_server/{server_id} [delete]
func DeleteServer(c *gin.Context) {
	serverID := c.Param("server_id")
	server, ok := posgresql_query.GetServerById(serverID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	status := posgresql_query.DeleteServer(serverID)
	switch status {
	case http.StatusOK:
		c.JSON(http.StatusOK, gin.H{
			"message":     "Server deleted successfully",
			"server_id":   server.ServerId,
			"server_name": server.ServerName,
			"server_ipv4": server.IPv4,
		})
	default:
		c.JSON(status, gin.H{"error": "Failed to delete server"})
		return
	}
}
