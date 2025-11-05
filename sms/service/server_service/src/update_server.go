package src

import (
	"net/http"
	"server_service/entities"
	elastic_query "server_service/infrastructure/elasticsearch/query"
	"time"

	"github.com/gin-gonic/gin"
)

// @Tags         Server
// @Summary      Update an existing server
// @Description  Update an existing server by validating input and updating server information
// @Accept       json
// @Produce      json
// @Param        jwt header string true "JWT token for authentication"
// @Param        request body entities.UpdateServerRequest true "Update server request"
// @Success      200 {object} entities.UpdateServerSuccessResponse "Server updated"
// @Failure      401 {object} entities.AuthErrorResponse "Authentication failed"
// @Failure      400 {object} entities.UpdateServerBadRequestResponse "Invalid request body"
// @Failure      404 {object} entities.UpdateServerStatusNotFoundResponse "Server not found"
// @Failure      409 {object} entities.UpdateServerConflictResponse "Server IP already exists"
// @Failure      429 {object} entities.RateLimitExceededResponse "Too many requests"
// @Failure      500 {object} entities.UpdateServerInternalServerErrorResponse "Internal server error"
// @Router       /servers/update-server [put]
func UpdateServer(c *gin.Context) {
	var req entities.UpdateServerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ServerId is required"})
		return
	}

	server, exists := elastic_query.GetServerById(req.Id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	if req.ServerName != "" {
		server.ServerName = req.ServerName
	}

	server.Status = req.Status
	server.LastUpdatedTime = time.Now().Unix()

	status := elastic_query.UpdateServerInfo(server)
	if status == http.StatusOK {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server updated successfully",
			"server":  server,
		})
	} else {
		c.JSON(status, gin.H{"error": "Failed to update server"})
	}
}
