package src

import (
	"net/http"

	"server_service/entities"
	elastic_query "server_service/infrastructure/elasticsearch/query"

	"github.com/gin-gonic/gin"
)

// @Tags         Server
// @Summary      Delete a server by ID
// @Description  Delete a server by its unique ID
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        request body entities.DeleteServerRequest true "Delete server request"
// @Success      200 {object} entities.DeleteServerResponse "Server deleted successfully"
// @Failure      401 {object} entities.AuthErrorResponse "Authentication failed"
// @Failure      404 {object} entities.DeleteServerStatusNotFoundResponse "Server not found"
// @Failure      500 {object} entities.DeleteServerInternalServerErrorResponse "Internal server error"
// @Router	     /server/delete_server [delete]
func DeleteServer(c *gin.Context) {
	var req entities.DeleteServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	server, ok := elastic_query.GetServerById(req.Id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	status := elastic_query.DeleteServer(req.Id)
	switch status {
	case http.StatusOK:
		c.JSON(http.StatusOK, gin.H{
			"message":     "Server deleted successfully",
			"_id":         req.Id,
			"server_name": server.ServerName,
			"server_ipv4": server.IPv4,
		})
	default:
		c.JSON(status, gin.H{"error": "Failed to delete server"})
		return
	}
}
