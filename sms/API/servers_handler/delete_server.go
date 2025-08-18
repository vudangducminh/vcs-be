package servers_handler

import (
	"net/http"

	"sms/object"
	elastic_query "sms/server/database/elasticsearch/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Servers
// @Summary      Delete a server by ID
// @Description  Delete a server by its unique ID
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        request body object.DeleteServerRequest true "Delete server request"
// @Success      200 {object} object.DeleteServerResponse "Server deleted successfully"
// @Failure      401 {object} object.AuthErrorResponse "Authentication failed"
// @Failure      404 {object} object.DeleteServerStatusNotFoundResponse "Server not found"
// @Failure      500 {object} object.DeleteServerInternalServerErrorResponse "Internal server error"
// @Router	     /servers/delete_server [delete]
func DeleteServer(c *gin.Context) {
	var req object.DeleteServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	server, ok := elastic_query.GetServerById(req.ServerId)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	status := elastic_query.DeleteServer(req.ServerId)
	switch status {
	case http.StatusOK:
		c.JSON(http.StatusOK, gin.H{
			"message":     "Server deleted successfully",
			"server_id":   req.ServerId,
			"server_name": server.ServerName,
			"server_ipv4": server.IPv4,
		})
	default:
		c.JSON(status, gin.H{"error": "Failed to delete server"})
		return
	}
}
